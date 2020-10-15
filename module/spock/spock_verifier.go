package spock

import (
	"fmt"
	"sync"

	"github.com/onflow/flow-go/crypto"
	"github.com/onflow/flow-go/engine"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/state/protocol"
)

// Verifier provides functionality to verify SPoCK proofs
type Verifier struct {
	sync.Mutex

	// state is used to query identities at a blockId to get StakingPublicKey
	protocolState protocol.ReadOnlyState

	// map of receipts by result ID that do not have matching SPoCK proofs
	// For instance, if there are 5 receipts that have 5 SpockSets, say, Spockset1, Spockset2, Spockset3, Spockset4, and Spockset5.
	// If Spockset1 and Spockset2 are for ER1, and they match with each other;
	// And Spockset3, Spockset4, Spockset5 are for ER2. and they don't match with each other.
	// Then there are 2 category, the first category has bucket: ER1_Spockset1, and the second category has buckets: ER2_Spockset3, ER2_Spockset4, and ER2_Spockset5.
	// When we receive each approval, we will first check which category it should go, and find all the buckets to try matching against.
	receipts map[flow.Identifier][]*flow.ExecutionReceipt

	// map to cache and keep track of matched approvals
	// approvals[resultID][approvalID] = receipt
	approvals map[flow.Identifier]map[flow.Identifier]*flow.ExecutionReceipt
}

// NewVerifier creates a new spock verifier
func NewVerifier(state protocol.ReadOnlyState) *Verifier {
	return &Verifier{
		protocolState: state,
		receipts:      make(map[flow.Identifier][]*flow.ExecutionReceipt),
		approvals:     make(map[flow.Identifier]map[flow.Identifier]*flow.ExecutionReceipt),
	}
}

// AddReceipt adds a receipt into map if the SPoCK proofs do not match any other receipt
func (v *Verifier) AddReceipt(receipt *flow.ExecutionReceipt) error {
	resultID := receipt.ExecutionResult.ID()

	// if receipts result id does not exist in map, create an array with the receipt
	_, ok := v.receipts[resultID]
	if !ok {
		receipts := make([]*flow.ExecutionReceipt, 0)
		receipts = append(receipts, receipt)
		v.receipts[resultID] = receipts
		return nil
	}

	// check if candidate receipt SPoCK proofs match proofs in map
	matched, err := v.matchReceipt(receipt)
	if err != nil {
		return fmt.Errorf("could not match receipt: %w", err)
	}

	// if not matched add to receipt array
	if !matched {
		v.receipts[resultID] = append(v.receipts[resultID], receipt)
	}

	// if matched we do nothing (strong-transitive property of the SPoCK scheme)

	return nil
}

// VerifyApproval verifies an approval with all the distinct receipts for the approvals
// result id and returns true if spocks match else false
func (v *Verifier) VerifyApproval(approval *flow.ResultApproval) (*flow.ExecutionReceipt, error) {

	// check if verification already done and that a match is there
	// note: we do not add to cache if the approval was not verified (sanity-check)
	receipt, ok := v.approvals[approval.Body.ExecutionResultID][approval.ID()]
	if receipt != nil && ok {
		return receipt, nil
	}

	// find identities
	approver, err := v.protocolState.AtBlockID(approval.Body.BlockID).Identity(approval.Body.ApproverID)
	if err != nil {
		return nil, fmt.Errorf("could not find approver identity")
	}

	receipts := v.receipts[approval.Body.ExecutionResultID]
	for _, receipt := range receipts {
		executor, err := v.protocolState.AtBlockID(receipt.ExecutionResult.BlockID).Identity(receipt.ExecutorID)
		if err != nil {
			return nil, fmt.Errorf("could not find executor identity")
		}

		// verify spock
		verified, err := crypto.SPOCKVerify(approver.StakingPubKey, approval.Body.Spock, executor.StakingPubKey, receipt.Spocks[approval.Body.ChunkIndex])
		if err != nil {
			return nil, fmt.Errorf("could not verify spocks: %w", err)
		}
		if verified {
			// map doesnt exist create it
			if _, ok = v.approvals[approval.Body.ExecutionResultID]; !ok {
				v.approvals[approval.Body.ExecutionResultID] = make(map[flow.Identifier]*flow.ExecutionReceipt)
			}

			// cache result for future verify calls for this approval
			v.approvals[approval.Body.ExecutionResultID][approval.ID()] = receipt

			return receipt, nil
		}
	}

	return nil, nil
}

// ClearReceipts clears all receipts for a specific resultID
func (v *Verifier) ClearReceipts(resultID flow.Identifier) bool {
	v.Lock()
	defer v.Unlock()

	// check if entry exists
	_, ok := v.receipts[resultID]
	if !ok {
		return false
	}
	_, ok = v.approvals[resultID]
	if !ok {
		return false
	}

	// clear receipts and approval cache
	delete(v.receipts, resultID)
	delete(v.approvals, resultID)

	return true
}

func (v *Verifier) matchReceipt(receipt *flow.ExecutionReceipt) (bool, error) {
	unmatchedReceipts := v.receipts[receipt.ExecutionResult.ID()]

	// get idenitity of candidate receipt
	identity, err := v.protocolState.AtBlockID(receipt.ExecutionResult.BlockID).Identity(receipt.ExecutorID)
	if err != nil {
		if protocol.IsIdentityNotFound(err) {
			return false, engine.NewInvalidInputErrorf("could not get executor identity: %w", err)
		}
		// unknown exception
		return false, fmt.Errorf("could not get executor identity: %w", err)
	}

	// for through each of the receipts to check for possible match of spocks
	// all the spocks in a receipt will have to match in order to be counted as
	// matched
MatchingReceiptsLoop:
	for _, u := range unmatchedReceipts {
		// we check if the receipt chunks length matches, if they dont match we return
		// false by default
		if u.ExecutionResult.Chunks.Len() != receipt.ExecutionResult.Chunks.Len() {
			return false, nil
		}

		// get receipt identity to get public key
		uIdentity, err := v.protocolState.AtBlockID(u.ExecutionResult.BlockID).Identity(u.ExecutorID)
		if err != nil {
			if protocol.IsIdentityNotFound(err) {
				return false, engine.NewInvalidInputErrorf("could not get executor identity: %w", err)
			}
			// unknown exception
			return false, fmt.Errorf("could not get executor identity: %w", err)
		}

		// attempt to match every spock in the receipt with the candidate receipt
		// if not verified then skip receipt
		for _, chunk := range u.ExecutionResult.Chunks {
			// check if spocks match
			verified, err := crypto.SPOCKVerify(identity.StakingPubKey, receipt.Spocks[chunk.Index], uIdentity.StakingPubKey, u.Spocks[chunk.Index])
			if err != nil {
				return false, fmt.Errorf("could not verify spocks: %w", err)
			}
			if !verified {
				continue MatchingReceiptsLoop
			}
		}

		// all spocks matched so we should exit for loop
		// since all spocks match we dont need to add this into an array
		return true, nil
	}

	return false, nil
}
package leader

import (
	"fmt"

	"github.com/onflow/flow-go/consensus/hotstuff/committee"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/flow/filter"
	"github.com/onflow/flow-go/model/indices"
	"github.com/onflow/flow-go/state/protocol"
)

// SelectionForEpoch pre-computes and returns leaders for the consensus committee
// in the given epoch.
// TODO: this should replace the consensus.go in this package
func SelectionForEpoch(epoch protocol.Epoch) (*committee.LeaderSelection, error) {

	// pre-compute leader selection for the current epoch
	identities, err := epoch.InitialIdentities()
	if err != nil {
		return nil, fmt.Errorf("could not get current epoch initial identities: %w", err)
	}
	seed, err := epoch.Seed(indices.ProtocolConsensusLeaderSelection...)
	if err != nil {
		return nil, fmt.Errorf("could not get current epoch seed: %w", err)
	}
	firstView, err := epoch.FirstView()
	if err != nil {
		return nil, fmt.Errorf("could not get current epoch first view: %w", err)
	}
	finalView, err := epoch.FinalView()
	if err != nil {
		return nil, fmt.Errorf("could not get current epoch final view: %w", err)
	}
	leaders, err := committee.ComputeLeaderSelectionFromSeed(
		firstView,
		seed,
		int(finalView-firstView+1), // add 1 because both first/final view are inclusive
		identities.Filter(filter.HasRole(flow.RoleConsensus)),
	)
	return leaders, err
}

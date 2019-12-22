// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package json

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/dapperlabs/flow-go/model/coldstuff"
	"github.com/dapperlabs/flow-go/model/flow"
	"github.com/dapperlabs/flow-go/model/trickle"
)

func encode(v interface{}) (*Envelope, error) {

	// determine the message type
	var code uint8
	switch v.(type) {

	case *trickle.Ping:
		code = CodePing
	case *trickle.Pong:
		code = CodePong
	case *trickle.Auth:
		code = CodeAuth
	case *trickle.Announce:
		code = CodeAnnounce
	case *trickle.Request:
		code = CodeRequest
	case *trickle.Response:
		code = CodeResponse

	case *flow.GuaranteedCollection:
		code = CodeGuaranteedCollection

	case *coldstuff.BlockProposal:
		code = CodeBlockProposal
	case *coldstuff.BlockVote:
		code = CodeBlockVote
	case *coldstuff.BlockCommit:
		code = CodeBlockCommit

	default:
		return nil, errors.Errorf("invalid encode type (%T)", v)
	}

	// encode the payload
	data, err := json.Marshal(v)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode payload")
	}

	env := Envelope{
		Code: code,
		Data: data,
	}

	return &env, nil
}

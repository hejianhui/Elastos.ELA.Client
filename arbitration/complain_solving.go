package arbitration

import (
	"ELAClient/crypto"
	"ELAClient/common"
)

const (
	Solving = iota
	Rejected
	Done
)

type ComplainListener interface {

	OnComplainFeedback([]byte)
}

type ComplainSolving interface {

	AcceptComplain(userKey *crypto.PublicKey, transactionHash common.Uint256) ([]byte, error)
	BroadcastComplainSolving([]byte) error

	GetComplainStatus(userKey *crypto.PublicKey, transactionHash common.Uint256) uint

	AddListener(listener ComplainListener) error
	RemoveListener(listener ComplainListener) error
}
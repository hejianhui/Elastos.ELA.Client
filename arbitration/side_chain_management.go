package arbitration

import (
	"ELAClient/crypto"
	"ELAClient/common"
)

type SideChain interface {
	AccountListener
	SpvValidation

	GetKey() *crypto.PublicKey
	GetNode() SideChainNode
	CreateDepositTransaction(target *crypto.PublicKey, information *SpvInformation) *TransactionInfo

	parseUserMainPublicKey(uint256 *common.Uint256) *crypto.PublicKey
}

type SideChainManager interface {

	Add(chain SideChain) error
	Remove(key *crypto.PublicKey) error

	GetChain(key *crypto.PublicKey) (SideChain, error)
	GetAllChains() ([]SideChain, error)
}
package arbitration

import "ELAClient/crypto"

type SideChain interface {
	AccountListener

	GetKey() *crypto.PublicKey
	GetNode() SideChainNode
	CreateDepositTransaction(target *crypto.PublicKey, information *SpvInformation) *TransactionInfo
}

type SideChainManager interface {

	Add(chain SideChain) error
	Remove(key *crypto.PublicKey) error

	GetChain(key *crypto.PublicKey) (SideChain, error)
	GetAllChains() ([]SideChain, error)
}
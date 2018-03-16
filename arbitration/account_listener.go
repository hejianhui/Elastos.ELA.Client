package arbitration

import (
	. "ELAClient/common"
	"ELAClient/crypto"
)

type AccountListener interface {

	OnUTXOChanged(transactionHash *Uint256) error
}

type AccountMonitor interface {

	SetAccount(account *crypto.PublicKey) error
	AddAccountListener(listener AccountListener) error

	fireUTXOChanged() error
}
package arbitration

import (
	"ELAClient/crypto"
	"ELAClient/common"
)

type MainChain interface {
	AccountListener
	SpvValidation

	CreateWithdrawTransaction()

	parseSideChainKey() *crypto.PublicKey
	parseUserSidePublicKey(uint256 *common.Uint256) *crypto.PublicKey
}
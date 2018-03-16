package arbitration

import (
	"ELAClient/crypto"
	"ELAClient/common"
)

type MainChain interface {
	AccountListener
	SpvValidation

	CreateWithdrawTransaction(withdrawBank *crypto.PublicKey, target *crypto.PublicKey) *TransactionInfo

	parseSideChainKey(uint256 *common.Uint256) *crypto.PublicKey
	parseUserSidePublicKey(uint256 *common.Uint256) *crypto.PublicKey
}
package arbitration

import (
	"ELAClient/wallet"
	"ELAClient/rpc"
	"ELAClient/common"
	"ELAClient/crypto"
)

func main() {

	// 流程前仲裁初始化
	var pkS *crypto.PublicKey
	var arbitratorGroup ArbitratorGroup
	currentArbitrator := arbitratorGroup.GetCurrentArbitrator()
	var mainAccountMonitor AccountMonitor
	mainAccountMonitor.SetAccount(pkS)
	mainAccountMonitor.AddAccountListener(currentArbitrator)

	//1. 钱包端
	var walletA wallet.Wallet
	var amount, fee *common.Fixed64
	var strAddressA, strAddressS string
	tx1, err := walletA.CreateTransaction(strAddressA, strAddressS, amount, fee)
	if tx1 == nil || err == nil {
		return
	}
	//sign tx1
	var transactionContent string
	rpc.CallAndUnmarshal("sendrawtransaction", rpc.Param("Data", transactionContent))

	//2. 仲裁主链
	//MainChain.OnUTXOChanged中逻辑（监听到充值申请）
	var transactionHash *common.Uint256
	pka := currentArbitrator.parseUserSidePublicKey(transactionHash)
	pkS = currentArbitrator.parseSideChainKey(transactionHash)
	spvInformation := currentArbitrator.GenerateSpvInformation(transactionHash)
	if valid, err := currentArbitrator.IsValid(spvInformation); !valid || err != nil {
		return
	}

	//3. 仲裁侧链
	sideChain, err := currentArbitrator.GetChain(pkS)
	tx2 := sideChain.CreateDepositTransaction(pka, spvInformation)
	sideChain.GetNode().SendTransaction(tx2)
}

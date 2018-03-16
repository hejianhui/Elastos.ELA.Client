package arbitration

import (
	"ELAClient/crypto"
	"ELAClient/wallet"
	"ELAClient/common"
	"ELAClient/rpc"
)

func main() {

	// 流程前仲裁初始化
	var pkDestroy *crypto.PublicKey
	var arbitratorGroup ArbitratorGroup
	currentArbitrator := arbitratorGroup.GetCurrentArbitrator()
	var sideAccountMonitor AccountMonitor
	sideAccountMonitor.SetAccount(pkDestroy)
	sideAccountMonitor.AddAccountListener(currentArbitrator)

	//1. 钱包端
	var walleta wallet.Wallet
	var amount, fee *common.Fixed64
	var strAddress_a, strAddressS string
	tx3, err := walleta.CreateTransaction(strAddress_a, strAddressS, amount, fee)
	if tx3 == nil || err == nil {
		return
	}
	//sign tx3
	var transactionContent string
	rpc.CallAndUnmarshal("sendrawtransaction", rpc.Param("Data", transactionContent))

	//2. 仲裁侧链
	//SideChain.OnUTXOChanged中逻辑（监听到提现申请）
	var transactionHash *common.Uint256
	sideChain, err := currentArbitrator.GetChain(pkDestroy)
	pkS := sideChain.GetKey()
	pkA := sideChain.parseUserMainPublicKey(transactionHash)
	spvInfo := sideChain.GenerateSpvInformation(transactionHash)
	if valid, err := sideChain.IsValid(spvInfo); !valid || err != nil {
		return
	}

	//3. 仲裁主链
	currentArbitrator.GetArbitrationNet().AddListener(currentArbitrator)
	tx4 := currentArbitrator.CreateWithdrawTransaction(pkS, pkA)
	tx4Bytes, err := tx4.Serialize()
	if err != nil {
		currentArbitrator.GetArbitrationNet().Broadcast(tx4Bytes)
	}

	//Arbitrator.OnReceived中逻辑（监听到其他仲裁人反馈）
	tx4.Deserialize(tx4Bytes)
	var tx4SignedContent string
	rpc.CallAndUnmarshal("sendrawtransaction", rpc.Param("Data", tx4SignedContent))
}

package arbitration

import (
	"ELAClient/crypto"
	"ELAClient/common"
	"fmt"
)

func main() {

	// 流程前仲裁初始化
	var arbitratorGroup ArbitratorGroup
	currentArbitrator := arbitratorGroup.GetCurrentArbitrator()
	currentArbitrator.GetComplainSolving().AddListener(currentArbitrator)

	// 1. 网页端发起申诉
	var userKey *crypto.PublicKey
	var transactionHash common.Uint256
	//send to current arbitrator

	// 2. 仲裁人
	solvingContent, err := currentArbitrator.GetComplainSolving().AcceptComplain(userKey, transactionHash)
	if err != nil {
		currentArbitrator.GetComplainSolving().BroadcastComplainSolving(solvingContent)
	}

	//Arbitrator.OnComplainFeedback中逻辑（监听其他仲裁人反馈，收集阶段完成）
	status := currentArbitrator.GetComplainSolving().GetComplainStatus(userKey, transactionHash)
	if status == Done {
		fmt.Println("Complain has been solved.")
	} else if status == Rejected {
		fmt.Println("Complain has been rejected.")
	}
}
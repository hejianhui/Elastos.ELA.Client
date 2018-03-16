package arbitration

type ArbitratorMain interface {
	MainChain
}

type ArbitratorSide interface {
	SideChainManager
}

type Arbitrator interface {
	ArbitratorMain
	ArbitratorSide

	GetArbitratorGroup() ArbitratorGroup
}
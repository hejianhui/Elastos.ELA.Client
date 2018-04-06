package wallet

import (
	"fmt"

	tx "github.com/elastos/Elastos.ELA.Client/core/transaction"
	. "github.com/elastos/Elastos.ELA.Client/rpc"
	. "github.com/elastos/Elastos.ELA.Utility/common"
	uti_tx "github.com/elastos/Elastos.ELA.Utility/core/transaction"
)

type DataSync interface {
	SyncChainData()
}

type DataSyncImpl struct {
	DataStore
	addresses []*Address
}

func GetDataSync(dataStore DataStore) DataSync {
	return &DataSyncImpl{
		DataStore: dataStore,
	}
}

func (sync *DataSyncImpl) SyncChainData() {
	// Get the addresses in this wallet
	sync.addresses, _ = sync.GetAddresses()

	var chainHeight uint32
	var currentHeight uint32
	var needSync bool

	for {
		chainHeight, currentHeight, needSync = sync.needSyncBlocks()
		if !needSync {
			break
		}

		for currentHeight <= chainHeight {
			block, err := GetBlockByHeight(currentHeight)
			if err != nil {
				break
			}
			sync.processBlock(block)

			// Update wallet height
			currentHeight = sync.CurrentHeight(block.BlockData.Height + 1)

			fmt.Print(">")
		}
	}

	fmt.Print("\n")
}

func (sync *DataSyncImpl) needSyncBlocks() (uint32, uint32, bool) {

	chainHeight, err := GetCurrentHeight()
	if err != nil {
		return 0, 0, false
	}

	currentHeight := sync.CurrentHeight(QueryHeightCode)

	if currentHeight >= chainHeight+1 {
		return chainHeight, currentHeight, false
	}

	return chainHeight, currentHeight, true
}

func (sync *DataSyncImpl) containAddress(address string) (*Address, bool) {
	for _, addr := range sync.addresses {
		if addr.Address == address {
			return addr, true
		}
	}
	return nil, false
}

func (sync *DataSyncImpl) processBlock(block *BlockInfo) {
	// Add UTXO to wallet address from transaction outputs
	for _, txn := range block.Transactions {

		// Add UTXOs to wallet address from transaction outputs
		for index, output := range txn.Outputs {
			if addr, ok := sync.containAddress(output.Address); ok {
				// Create UTXO input from output
				txHashBytes, _ := HexStringToBytesReverse(txn.Hash)
				referTxHash, _ := Uint256FromBytes(txHashBytes)
				lockTime := output.OutputLock
				if txn.TxType == uti_tx.CoinBase {
					lockTime = block.BlockData.Height + 100
				}
				amount, _ := StringToFixed64(output.Value)
				// Save UTXO input to data store
				addressUTXO := &AddressUTXO{
					Op:       tx.NewOutPoint(*referTxHash, uint16(index)),
					Amount:   amount,
					LockTime: lockTime,
				}
				sync.AddAddressUTXO(addr.ProgramHash, addressUTXO)
			}
		}

		// Delete UTXOs from wallet by transaction inputs
		for _, input := range txn.UTXOInputs {
			txHashBytes, _ := HexStringToBytesReverse(input.ReferTxID)
			referTxID, _ := Uint256FromBytes(txHashBytes)
			sync.DeleteUTXO(tx.NewOutPoint(*referTxID, input.ReferTxOutputIndex))
		}
	}
}

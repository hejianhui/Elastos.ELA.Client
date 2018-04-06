package transaction

import (
	"errors"

	"github.com/elastos/Elastos.ELA.Client/core/transaction/payload"
	. "github.com/elastos/Elastos.ELA.Utility/core/transaction"
)

const (
	TransferCrossChainAsset TransactionType = 0x08
)

type PayloadFactoryClientImpl struct {
	innerFactory *PayloadFactoryImpl
}

func (factor *PayloadFactoryClientImpl) Name(txType TransactionType) string {
	if name := factor.innerFactory.Name(txType); name != "Unknown" {
		return name
	}

	switch txType {
	case TransferCrossChainAsset:
		return "TransferCrossChainAsset"
	default:
		return "Unknown"
	}
}

func (factor *PayloadFactoryClientImpl) Create(txType TransactionType) (Payload, error) {
	if p, _ := factor.innerFactory.Create(txType); p != nil {
		return p, nil
	}

	switch txType {
	case TransferCrossChainAsset:
		return new(payload.TransferCrossChainAsset), nil
	default:
		return nil, errors.New("[NodeTransaction], invalid transaction type.")
	}
}

func init() {
	PayloadFactorySingleton = &PayloadFactoryClientImpl{innerFactory: &PayloadFactoryImpl{}}
}

package wallet

import (
	"errors"
	"strconv"

	"github.com/urfave/cli"

	"github.com/elastos/Elastos.ELA.Client/common/config"
	walt "github.com/elastos/Elastos.ELA.Client/wallet"
	. "github.com/elastos/Elastos.ELA.Utility/common"
	tx "github.com/elastos/Elastos.ELA.Utility/core/transaction"
)

var SingleTransactionCreatorSingleton SingleTransactionCreator

type SingleTransactionCreator interface {
	Create(c *cli.Context, param *SingleTransactionParameter, wallet walt.Wallet) (*tx.Transaction, error)
}

type SingleTransactionParameter struct {
	From   string
	Amount *Fixed64
	Fee    *Fixed64
}

type SingleTransactionCreatorImpl struct {
}

func (impl *SingleTransactionCreatorImpl) Create(c *cli.Context, param *SingleTransactionParameter, wallet walt.Wallet) (*tx.Transaction, error) {
	var to string
	standard := c.String("to")
	if standard != "" {
		to = standard
		lockStr := c.String("lock")
		if lockStr == "" {
			return wallet.CreateTransaction(param.From, to, param.Amount, param.Fee)
		} else {
			lock, err := strconv.ParseUint(lockStr, 10, 32)
			if err != nil {
				return nil, errors.New("invalid lock height")
			}
			return wallet.CreateLockedTransaction(param.From, to, param.Amount, param.Fee, uint32(lock))
		}
	} else {
		deposit := c.String("deposit")
		if deposit != "" {
			to := config.Params().DepositAddress
			return wallet.CreateCrossChainTransaction(param.From, to, deposit, param.Amount, param.Fee)
		}

		return nil, errors.New("use --to or --deposit or --withdraw to specify receiver address")
	}
}

func init() {
	SingleTransactionCreatorSingleton = &SingleTransactionCreatorImpl{}
}

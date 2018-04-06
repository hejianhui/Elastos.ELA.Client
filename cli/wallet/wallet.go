package wallet

import (
	"errors"
	"fmt"
	"os"

	"github.com/elastos/Elastos.ELA.Client/common/log"
	"github.com/elastos/Elastos.ELA.Client/common/password"
	"github.com/elastos/Elastos.ELA.Client/wallet"
	. "github.com/elastos/Elastos.ELA.Utility/common"

	"github.com/urfave/cli"
)

const (
	MinMultiSignKeys = 3
)

func printLine() {
	for i := 0; i < 80; i++ {
		fmt.Print("-")
	}
	fmt.Println()
}

func createWallet(name string, password []byte) error {

	password = getPassword(password, true)
	defer ClearBytes(password, len(password))

	_, err := wallet.Create(name, password)
	if err != nil {
		return err
	}
	return showAccountInfo(name, password)
}

func changePassword(name string, password []byte, wallet wallet.Wallet) error {
	// Verify old password
	oldPassword := getPassword(password, false)
	err := wallet.OpenKeystore(name, oldPassword)
	if err != nil {
		return err
	}
	defer ClearBytes(oldPassword, len(oldPassword))

	// Input new password
	fmt.Println("# input new password #")
	newPassword := getPassword(nil, true)
	if err := wallet.ChangePassword(oldPassword, newPassword); err != nil {
		return errors.New("failed to change password")
	}
	defer ClearBytes(newPassword, len(newPassword))

	fmt.Println("password changed successful")

	return nil
}

func showAccountInfo(name string, password []byte) error {
	password = getPassword(password, false)
	defer ClearBytes(password, len(password))

	keyStore, err := wallet.OpenKeystore(name, password)
	if err != nil {
		return err
	}
	programHash := keyStore.GetProgramHash()
	address, _ := programHash.ToAddress()
	publicKey := keyStore.GetPublicKey()
	publicKeyBytes, _ := publicKey.EncodePoint(true)

	printLine()
	fmt.Println("Address:     ", address)
	fmt.Println("Public Key:  ", BytesToHexString(publicKeyBytes))
	fmt.Println("ProgramHash: ", BytesToHexString(programHash.ToArrayReverse()))
	printLine()

	return nil
}

func listBalanceInfo(wallet wallet.Wallet) error {
	wallet.SyncChainData()
	addresses, err := wallet.GetAddresses()
	if err != nil {
		log.Error("Get addresses error:", err)
		return errors.New("get wallet addresses failed")
	}
	printLine()
	for _, address := range addresses {
		balance := Fixed64(0)
		programHash := address.ProgramHash
		UTXOs, err := wallet.GetAddressUTXOs(programHash)
		if err != nil {
			return errors.New("get " + address.Address + " UTXOs failed")
		}
		for _, utxo := range UTXOs {
			balance += *utxo.Amount
		}
		fmt.Println("Address:     ", address.Address)
		fmt.Println("ProgramHash: ", BytesToHexString(address.ProgramHash.ToArrayReverse()))
		fmt.Println("Balance:     ", balance.String())

		printLine()
	}
	return nil
}

func getPassword(passwd []byte, confirmed bool) []byte {
	var tmp []byte
	var err error
	if len(passwd) > 0 {
		tmp = []byte(passwd)
	} else {
		if confirmed {
			tmp, err = password.GetConfirmedPassword()
		} else {
			tmp, err = password.GetPassword()
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return tmp
}

func walletAction(context *cli.Context) {
	if context.NumFlags() == 0 {
		cli.ShowSubcommandHelp(context)
		os.Exit(0)
	}
	name := context.String("name")
	pass := context.String("password")

	// create wallet
	if context.Bool("create") {
		if err := createWallet(name, []byte(pass)); err != nil {
			fmt.Println("error: create wallet failed, ", err)
			cli.ShowCommandHelpAndExit(context, "create", 1)
		}
		return
	}

	wallet, err := wallet.Open()
	if err != nil {
		fmt.Println("error: open wallet failed, ", err)
		os.Exit(2)
	}

	// show account info
	if context.Bool("account") {
		if err := showAccountInfo(name, []byte(pass)); err != nil {
			fmt.Println("error: show account info failed, ", err)
			cli.ShowCommandHelpAndExit(context, "account", 3)
		}
		return
	}

	// change password
	if context.Bool("changepassword") {
		if err := changePassword(name, []byte(pass), wallet); err != nil {
			fmt.Println("error: change password failed, ", err)
			cli.ShowCommandHelpAndExit(context, "changepassword", 4)
		}
		return
	}

	// add account
	if pubKeyStr := context.String("addaccount"); pubKeyStr != "" {
		if err := addAccount(wallet, pubKeyStr); err != nil {
			fmt.Println("error: add standard account failed, ", err)
			cli.ShowCommandHelpAndExit(context, "addaccount", 5)
		}
		return
	}

	// add multi sign account
	if pubKeysStr := context.String("addmultisignaccount"); pubKeysStr != "" {
		if err := addMultiSignAccount(context, wallet, pubKeysStr); err != nil {
			fmt.Println("error: add multi sign account failed, ", err)
			cli.ShowCommandHelpAndExit(context, "addmultisignaccount", 5)
		}
		return
	}

	// delete account
	if address := context.String("deleteaccount"); address != "" {
		if err := deleteAccount(wallet, address); err != nil {
			fmt.Println("error: delete account failed, ", err)
			cli.ShowCommandHelpAndExit(context, "deleteaccount", 5)
		}
		return
	}

	// show addresses balance in this wallet
	if context.Bool("balance") {
		if err := listBalanceInfo(wallet); err != nil {
			fmt.Println("error: list balance info failed, ", err)
			cli.ShowCommandHelpAndExit(context, "balance", 6)
		}
		return
	}

	// transaction actions
	if param := context.String("transaction"); param != "" {
		switch param {
		case "create":
			if err := createTransaction(context, wallet); err != nil {
				fmt.Println("error:", err)
				os.Exit(701)
			}
		case "sign":
			if err := signTransaction(name, []byte(pass), context, wallet); err != nil {
				fmt.Println("error:", err)
				os.Exit(702)
			}
		case "send":
			if err := sendTransaction(context); err != nil {
				fmt.Println("error:", err)
				os.Exit(703)
			}
		default:
			cli.ShowCommandHelpAndExit(context, "transaction", 700)
		}
		return
	}

	// reset wallet
	if context.Bool("reset") {
		if err := wallet.Reset(); err != nil {
			fmt.Println("error: reset wallet data store failed, ", err)
			cli.ShowCommandHelpAndExit(context, "reset", 8)
		}
		fmt.Println("wallet data store was reset successfully")
		return
	}
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:        "wallet",
		Usage:       "wallet operations",
		Description: "With ela-cli wallet, you can create an account, check account balance or build, sign and send transactions.",
		ArgsUsage:   "[args]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "password, p",
				Usage: "arguments to pass the password value",
			},
			cli.StringFlag{
				Name:  "name, n",
				Usage: "to specify the created keystore file name or the keystore file path to open",
				Value: wallet.DefaultKeystoreFile,
			},
			cli.BoolFlag{
				Name:  "create, c",
				Usage: "create wallet, this will generate a keystore file within you account information",
			},
			cli.BoolFlag{
				Name:  "account, a",
				Usage: "show account address, public key and program hash",
			},
			cli.BoolFlag{
				Name:  "changepassword",
				Usage: "change the password to access this wallet, must do not forget it",
			},
			cli.BoolFlag{
				Name:  "reset",
				Usage: "clear the UTXOs stored in the local database",
			},
			cli.StringFlag{
				Name:  "addaccount",
				Usage: "add a standard account with a public key",
			},
			cli.StringFlag{
				Name: "addmultisignaccount",
				Usage: "add a multi-sign account with multiple public keys\n" +
					"\tuse -m to specify how many signatures are needed to create a valid transaction\n" +
					"\tby default M is public keys / 2 + 1, witch means greater than half",
			},
			cli.IntFlag{
				Name:  "m",
				Usage: "the M value to specify how many signatures are needed to create a valid transaction",
				Value: 0,
			},
			cli.StringFlag{
				Name:  "deleteaccount",
				Usage: "delete an account from database using it's address",
			},
			cli.BoolFlag{
				Name:  "balance, b",
				Usage: "list account balances stored in this wallet",
			},
			cli.StringFlag{
				Name: "transaction, t",
				Usage: "use [create, sign, send], to create, sign or send a transaction\n" +
					"\tcreate:\n" +
					"\t\tuse --to --amount --fee [--lock], or --file --fee [--lock]\n" +
					"\t\tto create a standard transaction, or multi output transaction\n" +
					"\tsign, send:\n" +
					"\t\tuse --file or --hex to specify the transaction file path or content\n",
			},
			cli.StringFlag{
				Name:  "from",
				Usage: "the spend address of the transaction",
			},
			cli.StringFlag{
				Name:  "to",
				Usage: "the receive address of the transaction",
			},
			cli.StringFlag{
				Name:  "amount",
				Usage: "the transfer amount of the transaction",
			},
			cli.StringFlag{
				Name:  "fee",
				Usage: "the transfer fee of the transaction",
			},
			cli.StringFlag{
				Name:  "lock",
				Usage: "the lock time to specify when the received asset can be spent",
			},
			cli.StringFlag{
				Name:  "hex",
				Usage: "the transaction content in hex string format to be sign or send",
			},
			cli.StringFlag{
				Name: "file, f",
				Usage: "the file path to specify a CSV file path with [address,amount] format as multi output content,\n" +
					"\tor the transaction file path with the hex string content to be sign or send",
			},
			cli.StringFlag{
				Name:  "key",
				Usage: "the public key of target account",
			},
			cli.StringFlag{
				Name:  "deposit",
				Usage: "create deposit transaction",
			},
			cli.StringFlag{
				Name:  "withdraw",
				Usage: "create withdraw transaction",
			},
		},
		Action: walletAction,
		OnUsageError: func(c *cli.Context, err error, subCommand bool) error {
			return cli.NewExitError(err, 1)
		},
	}
}

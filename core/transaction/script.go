package transaction

import (
	"bytes"
	"errors"
	"sort"

	. "github.com/elastos/Elastos.ELA.Client/crypto"
	. "github.com/elastos/Elastos.ELA.Utility/core/signature"
	"github.com/elastos/Elastos.ELA.Utility/crypto"
)

type OpCode byte

func CreateStandardRedeemScript(publicKey *crypto.PubKey) ([]byte, error) {
	content, err := publicKey.EncodePoint(true)
	if err != nil {
		return nil, errors.New("create standard redeem script, encode public key failed")
	}
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(len(content)))
	buf.Write(content)
	buf.WriteByte(byte(STANDARD))

	return buf.Bytes(), nil
}

func CreateMultiSignRedeemScript(M int, publicKeys []*crypto.PubKey) ([]byte, error) {
	// Write M
	opCode := OpCode(byte(PUSH1) + byte(M) - 1)
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(opCode))

	//sort pubkey
	sort.Sort(PubKeySlice(publicKeys))

	// Write public keys
	for _, pubkey := range publicKeys {
		content, err := pubkey.EncodePoint(true)
		if err != nil {
			return nil, errors.New("create multi sign redeem script, encode public key failed")
		}
		buf.WriteByte(byte(len(content)))
		buf.Write(content)
	}

	// Write N
	N := len(publicKeys)
	opCode = OpCode(byte(PUSH1) + byte(N) - 1)
	buf.WriteByte(byte(opCode))
	buf.WriteByte(MULTISIG)

	return buf.Bytes(), nil
}

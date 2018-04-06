package crypto

import (
	"crypto/ecdsa"
	"math/big"

	. "github.com/elastos/Elastos.ELA.Utility/crypto"
)

func NewPubKey(priKey []byte) *PubKey {
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = algSet.Curve

	k := new(big.Int)
	k.SetBytes(priKey)
	privateKey.D = k

	privateKey.PublicKey.X, privateKey.PublicKey.Y = algSet.Curve.ScalarBaseMult(k.Bytes())

	pubKey := new(PubKey)
	pubKey.X = privateKey.PublicKey.X
	pubKey.Y = privateKey.PublicKey.Y
	return pubKey
}

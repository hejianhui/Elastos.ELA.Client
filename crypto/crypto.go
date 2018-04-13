package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"

	. "github.com/elastos/Elastos.ELA.Utility/crypto"
)

func GenKeyPair() ([]byte, *PubKey, error) {

	privateKey, err := ecdsa.GenerateKey(AlgSet.Curve, rand.Reader)
	if err != nil {
		return nil, nil, errors.New("Generate key pair error")
	}

	publicKey := new(PubKey)
	publicKey.X = new(big.Int).Set(privateKey.PublicKey.X)
	publicKey.Y = new(big.Int).Set(privateKey.PublicKey.Y)

	return privateKey.D.Bytes(), publicKey, nil
}

func Sign(priKey []byte, data []byte) ([]byte, error) {

	digest := sha256.Sum256(data)

	privateKey := new(ecdsa.PrivateKey)
	privateKey.Curve = AlgSet.Curve
	privateKey.D = big.NewInt(0)
	privateKey.D.SetBytes(priKey)

	r := big.NewInt(0)
	s := big.NewInt(0)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, digest[:])
	if err != nil {
		fmt.Printf("Sign error\n")
		return nil, err
	}

	signature := make([]byte, SIGNATURELEN)

	lenR := len(r.Bytes())
	lenS := len(s.Bytes())
	copy(signature[SIGNRLEN-lenR:], r.Bytes())
	copy(signature[SIGNATURELEN-lenS:], s.Bytes())
	return signature, nil
}

type PubKeySlice []*PubKey

func (p PubKeySlice) Len() int { return len(p) }
func (p PubKeySlice) Less(i, j int) bool {
	r := p[i].X.Cmp(p[j].X)
	if r <= 0 {
		return true
	}
	return false
}
func (p PubKeySlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func Equal(e1 *PubKey, e2 *PubKey) bool {
	r := e1.X.Cmp(e2.X)
	if r != 0 {
		return false
	}
	r = e1.Y.Cmp(e2.Y)
	if r == 0 {
		return true
	}
	return false
}

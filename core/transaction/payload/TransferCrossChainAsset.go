package payload

import (
	"ELAClient/common/serialization"
	"errors"
	"io"
)

type TransferCrossChainAsset struct {
	//string: targetAddress; uint64: output index
	CrossChainAddress map[string]uint64
}

func (a *TransferCrossChainAsset) Data(version byte) []byte {
	//TODO: implement TransferCrossChainAsset.Data()
	return []byte{0}
}

func (a *TransferCrossChainAsset) Serialize(w io.Writer, version byte) error {
	if a.CrossChainAddress == nil {
		return errors.New("Invalid publickey map")
	}

	if err := serialization.WriteVarUint(w, uint64(len(a.CrossChainAddress))); err != nil {
		return errors.New("publicKey map's length serialize failed")
	}

	for k, v := range a.CrossChainAddress {
		if err := serialization.WriteVarString(w, k); err != nil {
			return errors.New("publicKey map's key serialize failed")
		}

		if err := serialization.WriteVarUint(w, v); err != nil {
			return errors.New("publicKey map's value serialize failed")
		}
	}

	return nil
}

func (a *TransferCrossChainAsset) Deserialize(r io.Reader, version byte) error {
	if a.CrossChainAddress == nil {
		return errors.New("Invalid public key map")
	}

	length, err := serialization.ReadVarUint(r, 0)
	if err != nil {
		return errors.New("publicKey map's length deserialize failed")
	}

	a.CrossChainAddress = nil
	a.CrossChainAddress = make(map[string]uint64)
	for i := uint64(0); i < length; i++ {
		k, err := serialization.ReadVarString(r)
		if err != nil {
			return errors.New("publicKey map's key deserialize failed")
		}

		v, err := serialization.ReadVarUint(r, 0)
		if err != nil {
			return errors.New("publicKey map's value deserialize failed")
		}

		a.CrossChainAddress[k] = v
	}

	return nil
}

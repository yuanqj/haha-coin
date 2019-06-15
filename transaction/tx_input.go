package transaction

import (
	"bytes"
	"github.com/yuanqj/haha-coin/wallet"
)

type TXInput struct {
	TxID      *TxIDType
	OutIdx    int
	Signature []byte
	PubKey    []byte
}

func (in *TXInput) UsesKey(pubKeyHash []byte) (bool, error) {
	lockingHash, err := wallet.HashPubKey(in.PubKey)
	if err != nil {
		return false, err
	}
	return bytes.Compare(lockingHash, pubKeyHash) == 0, nil
}

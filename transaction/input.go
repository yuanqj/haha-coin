package transaction

import (
	"bytes"
	"github.com/yuanqj/haha-coin/wallet"
)

type Input struct {
	TxID      *IDType
	OutputIdx int
	Signature []byte
	PubKey    []byte
}

func (in *Input) UsesKey(pubKeyHash []byte) (bool, error) {
	lockingHash, err := wallet.HashPubKey(in.PubKey)
	if err != nil {
		return false, err
	}
	return bytes.Compare(lockingHash, pubKeyHash) == 0, nil
}

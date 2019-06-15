package transaction

import (
	"bytes"
	"github.com/yuanqj/haha-coin/base58"
)

type TXOutput struct {
	Val        int
	PubKeyHash []byte
}

type TXOutputKey struct {
	TxID TxIDType
	Idx  int
}

type TXOutputWraper struct {
	Key *TXOutputKey
	Out *TXOutput
}

func NewTXOutput(val int, addr string) (*TXOutput, error) {
	out := &TXOutput{val, nil}
	if err := out.Lock(addr); err != nil {
		return nil, err
	}
	return out, nil
}

func (out *TXOutput) Lock(addr string) error {
	pubKeyHash, err := base58.Decode(addr)
	if err != nil {
		return err
	}
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
	return nil
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

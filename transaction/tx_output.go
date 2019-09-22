package transaction

import (
	"bytes"
	"github.com/yuanqj/haha-coin/base58"
)

type Output struct {
	Val        int
	PubKeyHash []byte
}

type TXOutputKey struct {
	TxID IDType
	Idx  int
}

type TXOutputWraper struct {
	Key    *TXOutputKey
	Output *Output
}

func NewTXOutput(val int, addr string) (*Output, error) {
	out := &Output{val, nil}
	if err := out.Lock(addr); err != nil {
		return nil, err
	}
	return out, nil
}

func (out *Output) Lock(addr string) error {
	pubKeyHash, err := base58.Decode(addr)
	if err != nil {
		return err
	}
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
	return nil
}

func (out *Output) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

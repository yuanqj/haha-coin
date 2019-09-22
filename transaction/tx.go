package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/yuanqj/haha-coin/wallet"
	"math/big"
	"strings"
)

const subsidy = 10

type IDType [32]byte

type Transaction struct {
	ID      *IDType
	Inputs  []*Input
	Outputs []*Output
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs[0].TxID) == 0 && tx.Inputs[0].OutputIdx == -1
}

func (tx Transaction) Serialize() ([]byte, error) {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	if err := enc.Encode(tx); err != nil {
		return nil, err
	}
	return encoded.Bytes(), nil
}

func (tx *Transaction) Hash() (*IDType, error) {
	cont, err := tx.Trim().Serialize()
	if err != nil {
		return nil, err
	}
	hash := IDType(sha256.Sum256(cont))
	return &hash, nil
}

func NewCoinbaseTransaction(to string) (*Transaction, error) {
	txIn := &Input{OutputIdx: -1, Signature: nil, PubKey: []byte("Reward")}
	txOut, err := NewTXOutput(subsidy, to)
	if err != nil {
		return nil, err
	}
	tx := &Transaction{Inputs: []*Input{txIn}, Outputs: []*Output{txOut}}
	if id, err := tx.Hash(); err != nil {
		return nil, err
	} else {
		tx.ID = id
		return tx, nil
	}
}

func NewUTXOTransaction(srcWallet *wallet.Wallet, dstAddr string, amt int, utxos []*TXOutputWraper) (tx *Transaction, err error) {
	tot := 0
	for _, utxo := range utxos {
		tot += utxo.Output.Val
	}
	if tot < amt {
		err = fmt.Errorf("no enough blance")
		return
	}

	// Inputs
	ins := make([]*Input, len(utxos))
	for i, utxo := range utxos {
		ins[i] = &Input{TxID: &utxo.Key.TxID, OutputIdx: utxo.Key.Idx, PubKey: srcWallet.PubKey}
	}

	// Outputs
	outs := make([]*Output, 2)
	if outs[0], err = NewTXOutput(amt, dstAddr); err != nil {
		return
	}
	if left := tot - amt; left > 0 {
		if outs[1], err = NewTXOutput(left, srcWallet.Addr); err != nil {
			return
		}
	} else {
		outs = outs[:1]
	}

	tx = &Transaction{Outputs: outs, Inputs: ins}
	if tx.ID, err = tx.Hash(); err != nil {
		return
	}
	return
}

func (tx *Transaction) Trim() Transaction {
	ins := make([]*Input, len(tx.Inputs))
	outs := make([]*Output, len(tx.Outputs))
	for i, in := range tx.Inputs {
		ins[i] = &Input{in.TxID, in.OutputIdx, nil, nil}
	}
	for i, out := range tx.Outputs {
		outs[i] = &Output{out.Val, out.PubKeyHash}
	}
	return Transaction{tx.ID, ins, outs}
}

func (tx *Transaction) Sign(prv *ecdsa.PrivateKey, prevTXs []*Transaction) error {
	if tx.IsCoinbase() {
		return nil
	}
	for i := range tx.Inputs {
		if prevTXs[i].ID == nil {
			return fmt.Errorf("invalid previous trasaction")
		}
	}

	trimmed := tx.Trim()
	for i, in := range trimmed.Inputs {
		prevTX := prevTXs[i]
		in.PubKey = prevTX.Outputs[in.OutputIdx].PubKeyHash
		id, err := trimmed.Hash()
		if err != nil {
			return err
		}
		trimmed.ID = id
		in.PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, prv, id[:])
		if err != nil {
			return err
		}
		sig := append(r.Bytes(), s.Bytes()...)
		in.Signature = sig
	}
	return nil
}

func (tx *Transaction) Verify(prevTXs []*Transaction) (bool, error) {
	if tx.IsCoinbase() {
		return true, nil
	}
	for i := range tx.Inputs {
		if prevTXs[i].ID == nil {
			return false, fmt.Errorf("invalid previous trasaction")
		}
	}
	trimmed := tx.Trim()
	c := elliptic.P256()
	for i, in := range tx.Inputs {
		prevTX := prevTXs[i]
		trimmed.Inputs[i].Signature = nil
		trimmed.Inputs[i].PubKey = prevTX.Outputs[in.OutputIdx].PubKeyHash
		id, err := trimmed.Hash()
		if err != nil {
			return false, err
		}
		trimmed.ID = id
		trimmed.Inputs[i].PubKey = nil

		r, s := new(big.Int), new(big.Int)
		lenSig := len(in.Signature)
		r.SetBytes(in.Signature[:lenSig/2])
		s.SetBytes(in.Signature[lenSig/2:])

		x, y := new(big.Int), new(big.Int)
		lenKey := len(in.PubKey)
		x.SetBytes(in.PubKey[:lenKey/2])
		y.SetBytes(in.PubKey[lenKey/2:])
		pubKeyRaw := &ecdsa.PublicKey{Curve: c, X: x, Y: y}
		if !ecdsa.Verify(pubKeyRaw, trimmed.ID[:], r, s) {
			return false, nil
		}
	}
	return true, nil
}

func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))

	for i, in := range tx.Inputs {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", in.TxID))
		lines = append(lines, fmt.Sprintf("       Output:       %d", in.OutputIdx))
		lines = append(lines, fmt.Sprintf("       Signature: %x", in.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", in.PubKey))
	}

	for i, out := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", out.Val))
		lines = append(lines, fmt.Sprintf("       Script: %x", out.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

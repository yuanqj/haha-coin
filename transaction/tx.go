package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"haha/wallet"
	"math/big"
	"strings"
)

const subsidy = 10

type TxIDType [32]byte

type Transaction struct {
	ID   *TxIDType
	Ins  []*TXInput
	Outs []*TXOutput
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Ins[0].TxID) == 0 && tx.Ins[0].OutIdx == -1
}

func (tx Transaction) Serialize() ([]byte, error) {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	if err := enc.Encode(tx); err != nil {
		return nil, err
	}
	return encoded.Bytes(), nil
}

func (tx *Transaction) Hash() (*TxIDType, error) {
	cont, err := tx.Trim().Serialize()
	if err != nil {
		return nil, err
	}
	hash := TxIDType(sha256.Sum256(cont))
	return &hash, nil
}

func NewCoinbaseTransaction(to string) (*Transaction, error) {
	txIn := &TXInput{OutIdx: -1, Signature: nil, PubKey: []byte("Reward")}
	txOut, err := NewTXOutput(subsidy, to)
	if err != nil {
		return nil, err
	}
	tx := &Transaction{Ins: []*TXInput{txIn}, Outs: []*TXOutput{txOut}}
	if id, err := tx.Hash(); err != nil {
		return nil, err
	} else {
		tx.ID = id
		return tx, nil
	}
}

func NewUTXOTransaction(fromWallet *wallet.Wallet, toAddr string, amt int, utxos []*TXOutputWraper) (tx *Transaction, err error) {
	tot := 0
	for _, utxo := range utxos {
		tot += utxo.Out.Val
	}
	if tot < amt {
		err = fmt.Errorf("no enough blance")
		return
	}

	// Inputs
	ins := make([]*TXInput, len(utxos))
	for i, utxo := range utxos {
		ins[i] = &TXInput{TxID: &utxo.Key.TxID, OutIdx: utxo.Key.Idx, PubKey: fromWallet.PubKey}
	}

	// Outputs
	outs := make([]*TXOutput, 2)
	if outs[0], err = NewTXOutput(amt, toAddr); err != nil {
		return
	}
	if left := tot - amt; left > 0 {
		if outs[1], err = NewTXOutput(left, fromWallet.Addr); err != nil {
			return
		}
	} else {
		outs = outs[:1]
	}

	tx = &Transaction{Outs: outs, Ins: ins}
	if tx.ID, err = tx.Hash(); err != nil {
		return
	}
	return
}

func (tx *Transaction) Trim() Transaction {
	ins := make([]*TXInput, len(tx.Ins))
	outs := make([]*TXOutput, len(tx.Outs))
	for i, in := range tx.Ins {
		ins[i] = &TXInput{in.TxID, in.OutIdx, nil, nil}
	}
	for i, out := range tx.Outs {
		outs[i] = &TXOutput{out.Val, out.PubKeyHash}
	}
	return Transaction{tx.ID, ins, outs}
}

func (tx *Transaction) Sign(prv *ecdsa.PrivateKey, prevTXs []*Transaction) error {
	if tx.IsCoinbase() {
		return nil
	}
	for i := range tx.Ins {
		if prevTXs[i].ID == nil {
			return fmt.Errorf("invalid previous trasaction")
		}
	}

	trimmed := tx.Trim()
	for i, in := range trimmed.Ins {
		prevTX := prevTXs[i]
		in.PubKey = prevTX.Outs[in.OutIdx].PubKeyHash
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
	for i := range tx.Ins {
		if prevTXs[i].ID == nil {
			return false, fmt.Errorf("invalid previous trasaction")
		}
	}
	trimmed := tx.Trim()
	c := elliptic.P256()
	for i, in := range tx.Ins {
		prevTX := prevTXs[i]
		trimmed.Ins[i].Signature = nil
		trimmed.Ins[i].PubKey = prevTX.Outs[in.OutIdx].PubKeyHash
		id, err := trimmed.Hash()
		if err != nil {
			return false, err
		}
		trimmed.ID = id
		trimmed.Ins[i].PubKey = nil

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

	for i, in := range tx.Ins {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", in.TxID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", in.OutIdx))
		lines = append(lines, fmt.Sprintf("       Signature: %x", in.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", in.PubKey))
	}

	for i, out := range tx.Outs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", out.Val))
		lines = append(lines, fmt.Sprintf("       Script: %x", out.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

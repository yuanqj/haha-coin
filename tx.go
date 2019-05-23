package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

const subsidy = 10

type TxIDType [32]byte

type Transaction struct {
	ID   TxIDType
	Ins  []TXInput
	Outs []TXOutput
}

type TXOutput struct {
	Val          int
	ScriptPubKey string
}

type TXOutputKey struct {
	TxID TxIDType
	Idx  int
}

type TXOutputWraper struct {
	Key *TXOutputKey
	Out *TXOutput
}

type TXInput struct {
	TxID      TxIDType
	OutIdx    int
	ScriptSig string
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Ins[0].TxID) == 0 && tx.Ins[0].OutIdx == -1
}

func (tx *Transaction) SetID() error {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	if err := encoder.Encode(tx); err != nil {
		return err
	}
	hash := sha256.Sum256(encoded.Bytes())
	tx.ID = hash
	return nil
}

func NewCoinbaseTransaction(to string) (*Transaction, error) {
	txIn := TXInput{OutIdx: -1, ScriptSig: "Reward"}
	txOut := TXOutput{subsidy, to}
	tx := &Transaction{Ins: []TXInput{txIn}, Outs: []TXOutput{txOut}}
	if err := tx.SetID(); err != nil {
		return nil, err
	} else {
		return tx, nil
	}
}

func NewUTXOTransaction(from, to string, amt int, bc *Blockchain) (tx *Transaction, err error) {
	utxos, tot, err := bc.UTXOs(from, amt)
	if err != nil {
		return
	}
	if tot < amt {
		err = fmt.Errorf("no enough blance")
		return
	}

	// Inputs
	ins := make([]TXInput, len(utxos))
	for i, utxo := range utxos {
		ins[i].TxID = utxo.Key.TxID
		ins[i].OutIdx = utxo.Key.Idx
		ins[i].ScriptSig = from
	}

	// Outputs
	out := TXOutput{Val: amt, ScriptPubKey: to}
	outs := []TXOutput{out}
	if left := tot - amt; left > 0 {
		out := TXOutput{Val: left, ScriptPubKey: from}
		outs = append(outs, out)
	}

	tx = &Transaction{Outs: outs, Ins: ins}
	tx.SetID()
	return
}

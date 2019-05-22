package main

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp           int64
	PrevBlockHash, Hash []byte
	Transactions        []*Transaction
	Nonce               int
}

func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{Timestamp: time.Now().Unix(), PrevBlockHash: prevBlockHash, Transactions: txs}
	pow := NewPoW(block)
	nonce, hash := pow.Run()
	block.Hash, block.Nonce = hash[:], nonce
	return block
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func (b *Block) Serialize() ([]byte, error) {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	if err := encoder.Encode(b); err != nil {
		return nil, err
	} else {
		return buff.Bytes(), nil
	}

}

func DeserializeBlock(hex []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(hex))
	if err := decoder.Decode(&block); err != nil {
		return nil, err
	} else {
		return &block, nil
	}
}

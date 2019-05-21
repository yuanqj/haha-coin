package main

import (
	"time"
)

type Block struct {
	Timestamp           int64
	PrevBlockHash, Hash []byte
	Data                []byte
	Nonce               int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	b := &Block{Timestamp: time.Now().Unix(), PrevBlockHash: prevBlockHash, Data: []byte(data)}
	pow := NewPoW(b)
	nonce, hash := pow.Run()
	b.Hash, b.Nonce = hash[:], nonce
	return b
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

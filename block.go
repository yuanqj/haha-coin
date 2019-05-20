package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp           int64
	PrevBlockHash, Hash []byte
	Data                []byte
}

func (b *Block) SetHash() {
	ts := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, ts, b.Data}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	b := &Block{Timestamp: time.Now().Unix(), PrevBlockHash: prevBlockHash, Data: []byte(data)}
	b.SetHash()
	return b
}

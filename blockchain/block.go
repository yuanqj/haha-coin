package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
	"github.com/yuanqj/haha-coin/transaction"
)

type Block struct {
	Timestamp           int64
	PrevBlockHash, Hash []byte
	Transactions        []*transaction.Transaction
	Nonce               int
}

func NewBlock(txs []*transaction.Transaction, prevBlockHash []byte) *Block {
	block := &Block{Timestamp: time.Now().Unix(), PrevBlockHash: prevBlockHash, Transactions: txs}
	pow := NewPoW(block)
	nonce, hash := pow.Run()
	block.Hash, block.Nonce = hash[:], nonce
	return block
}

func NewGenesisBlock(coinbase *transaction.Transaction) *Block {
	return NewBlock([]*transaction.Transaction{coinbase}, []byte{})
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID[:])
	}
	txHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
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

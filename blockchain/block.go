package blockchain

import (
	"bytes"
	"encoding/gob"
	"github.com/yuanqj/haha-coin/merkle"
	"github.com/yuanqj/haha-coin/transaction"
	"time"
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
	block.Nonce, block.Hash = pow.Run()
	return block
}

func NewGenesisBlock(coinbase *transaction.Transaction) *Block {
	return NewBlock([]*transaction.Transaction{coinbase}, []byte{})
}

func (b *Block) HashTransactions() []byte {
	txHashes := make([][]byte, 0, len(b.Transactions))
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID[:])
	}
	tree := merkle.NewTree(txHashes)
	return tree.Root.Hash[:]
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

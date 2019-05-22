package main

import (
	bolt "github.com/etcd-io/bbolt"
)

const dbFile = "haha.db"

var (
	bucketBlocks = []byte("blocks")
	keyLastBlock = []byte("haha")
)

type Blockchain struct {
	db  *bolt.DB
	tip []byte
}

type BlockchainIterator struct {
	currHash []byte
	db       *bolt.DB
}

func NewBlockchain() (bc *Blockchain, err error) {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return
	}
	err = db.Update(
		func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(bucketBlocks))
			if bucket == nil {
				if bucket, err = tx.CreateBucket(bucketBlocks); err != nil {
					return err
				}

				genesis := NewGenesisBlock()
				encodedBlock, err := genesis.Serialize()
				if err != nil {
					return err
				}
				if err := bucket.Put(genesis.Hash, encodedBlock); err != nil {
					return err
				}
				if err := bucket.Put(keyLastBlock, genesis.Hash); err != nil {
					return err
				}
			}
			tip = bucket.Get(keyLastBlock)
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &Blockchain{db: db, tip: tip}, nil
}

func (bc *Blockchain) AddBlock(data string) (err error) {
	var lastHash []byte
	err = bc.db.View(
		func(tx *bolt.Tx) error {
			bucket := tx.Bucket(bucketBlocks)
			lastHash = bucket.Get(keyLastBlock)
			return nil
		},
	)
	if err != nil {
		return
	}

	block := NewBlock(data, lastHash)
	encodedBlock, err := block.Serialize()
	if err != nil {
		return
	}
	err = bc.db.Update(
		func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(bucketBlocks))
			if err := bucket.Put(block.Hash, encodedBlock); err != nil {
				return err
			}
			if err := bucket.Put(keyLastBlock, block.Hash); err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return
	}
	bc.tip = block.Hash
	return
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{currHash: bc.tip, db: bc.db}
}

func (bci *BlockchainIterator) Next() (*Block, error) {
	var block *Block
	err := bci.db.View(
		func(tx *bolt.Tx) error {
			bucket := tx.Bucket(bucketBlocks)
			encodedBlock := bucket.Get(bci.currHash)
			var err error
			block, err = DeserializeBlock(encodedBlock)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	bci.currHash = block.PrevBlockHash
	return block, nil
}

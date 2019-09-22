package blockchain

import (
	bolt "github.com/etcd-io/bbolt"
)

type Iterator struct {
	currHash []byte
	db       *bolt.DB
}

func (bci *Iterator) Next() (*Block, error) {
	if len(bci.currHash) <= 0 {
		return nil, nil
	}
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

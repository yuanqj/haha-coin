package main

import (
	"fmt"
	bolt "github.com/etcd-io/bbolt"
	"os"
)

const dbFile = "haha.db"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

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

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func LoadBlockchain() (bc *Blockchain, err error) {
	if !dbExists() {
		return nil, fmt.Errorf("no existing blockchain found")
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return
	}
	err = db.Update(
		func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(bucketBlocks))
			if bucket == nil {
				return fmt.Errorf("bucket of blocks not found")
			}
			tip = bucket.Get(keyLastBlock)
			return nil
		},
	)
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Blockchain{db: db, tip: tip}, nil
}

func CreateBlockchain(addr string) (bc *Blockchain, err error) {
	if dbExists() {
		return nil, fmt.Errorf("there already exists a blockchain")
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return
	}
	err = db.Update(
		func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucket(bucketBlocks)
			if err != nil {
				return err
			}
			txGenesis, err := NewCoinbaseTransaction(addr)
			if err != nil {
				return err
			}
			genesis := NewGenesisBlock(txGenesis)
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
			tip = genesis.Hash
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &Blockchain{db: db, tip: tip}, nil
}

func (bc *Blockchain) MineBlock(txs []*Transaction) (err error) {
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

	block := NewBlock(txs, lastHash)
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

func (bc *Blockchain) FindUTXOs(addr string, amt int) (utxos []*TXOutputWraper, tot int, err error) {
	stxos := make(map[TXOutputKey]bool)
	bci := bc.Iterator()

	Blocks:
	for {
		block, errBlock := bci.Next()
		if errBlock != nil {
			err = errBlock
			break
		}
		if block == nil {
			break
		}

		for _, tx := range block.Transactions {
			// Outputs
			for idx, out := range tx.Outs {
				key := TXOutputKey{TxID: tx.ID, Idx: idx}
				if !stxos[key] { // Unspent
					utxo := &TXOutputWraper{Key: &key, Out: &out}
					utxos = append(utxos, utxo)
					tot += out.Val
				}
				if tot >= amt {
					break Blocks
				}
			}

			// Inputs
			for _, in := range tx.Ins {
				key := TXOutputKey{TxID: tx.ID, Idx: in.Vout}
				stxos[key] = true // Spent
			}
		}
	}
	return
}

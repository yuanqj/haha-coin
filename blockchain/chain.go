package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	bolt "github.com/etcd-io/bbolt"
	"haha/transaction"
	"haha/wallet"
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

func (bc *Blockchain) Close() {
	bc.db.Close()
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
			txGenesis, err := transaction.NewCoinbaseTransaction(addr)
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

func (bc *Blockchain) MineBlock(txs []*transaction.Transaction) (err error) {
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

func (bc *Blockchain) UTXOs(w *wallet.Wallet, amt int) (utxos []*transaction.TXOutputWraper, tot int, err error) {
	stxos := make(map[transaction.TXOutputKey]bool)
	bci := bc.Iterator()

	pubKeyHash, err := wallet.HashPubKey(w.PubKey)
	if err != nil {
		return nil, 0, err
	}

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
				if !out.IsLockedWithKey(pubKeyHash) {
					continue
				}
				key := transaction.TXOutputKey{TxID: *tx.ID, Idx: idx}
				if !stxos[key] { // Unspent
					utxo := &transaction.TXOutputWraper{Key: &key, Out: out}
					utxos = append(utxos, utxo)
					tot += out.Val
				}
				if tot >= amt {
					break Blocks
				}
			}

			// Inputs
			if tx.IsCoinbase() {
				continue
			}
			for _, in := range tx.Ins {
				use, err := in.UsesKey(pubKeyHash)
				if err != nil {
					return nil, 0, err
				}
				if !use {
					continue
				}
				key := transaction.TXOutputKey{TxID: *in.TxID, Idx: in.OutIdx}
				stxos[key] = true // Spent
			}
		}
	}
	return
}

func (bc *Blockchain) FindTransaction(ID *transaction.TxIDType) (tx *transaction.Transaction, err error) {
	bci := bc.Iterator()
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
			if bytes.Compare(tx.ID[:], ID[:]) == 0 {
				return tx, nil
			}
		}
	}
	return nil, fmt.Errorf("transaction not found")
}

func (bc *Blockchain) SignTransaction(tx *transaction.Transaction, prv *ecdsa.PrivateKey) (err error) {
	prevTXs := make([]*transaction.Transaction, len(tx.Ins))
	for i, in := range tx.Ins {
		if prevTXs[i], err = bc.FindTransaction(in.TxID); err != nil {
			return
		}
	}
	tx.Sign(prv, prevTXs)
	return
}

func (bc *Blockchain) VerifyTransaction(tx *transaction.Transaction) (valid bool, err error) {
	prevTXs := make([]*transaction.Transaction, len(tx.Ins))
	for i, in := range tx.Ins {
		if prevTXs[i], err = bc.FindTransaction(in.TxID); err != nil {
			return
		}
	}
	return tx.Verify(prevTXs)
}

package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const MaxNonce = math.MaxInt64
const TargetBits uint = 20

type PoW struct {
	block  *Block
	target *big.Int
}

func NewPoW(block *Block) *PoW {
	one := big.NewInt(1)
	return &PoW{block: block, target: one.Lsh(one, 256-TargetBits)}
}

func (pow *PoW) prepareData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Hash,
			Int2Hex(pow.block.Timestamp),
			Int2Hex(int64(TargetBits)),
			Int2Hex(int64(nonce)),
		},
		[]byte{},
	)
}

func (pow *PoW) Run() (int, []byte) {
	fmt.Printf("\n>>>>>>> Mining...\n")
	fmt.Printf("# PrevHash: %x\n", pow.block.PrevBlockHash)

	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	for nonce < MaxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		}
		nonce++
	}
	fmt.Printf("# Hash: %x\n\n", hash)
	return nonce, hash[:]
}

func (pow *PoW) Validate() bool {
	data := pow.prepareData(pow.block.Nonce)
	var hashInt big.Int
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1
}

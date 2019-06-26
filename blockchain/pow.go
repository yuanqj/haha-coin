package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/yuanqj/haha-coin/util"
	"math"
	"math/big"
)

const MaxNonce = math.MaxInt64
const TargetBits uint = 20

var (
	one    = big.NewInt(1)
	target = one.Lsh(one, 256-TargetBits)
)

type PoW struct {
	block  *Block
	target *big.Int
}

func NewPoW(block *Block) *PoW {
	return &PoW{block: block, target: target}
}

func (pow *PoW) prepareData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			util.Int2Hex(pow.block.Timestamp),
			util.Int2Hex(int64(TargetBits)),
			util.Int2Hex(int64(nonce)),
		},
		[]byte{},
	)
}

func (pow *PoW) Run() (int, []byte) {
	fmt.Printf("\n>>>>>>> Mining...\n")
	fmt.Printf("# PrevHash: %x\n", pow.block.PrevBlockHash)

	var val big.Int
	var hash []byte
	nonce := 0
	for ; nonce < MaxNonce; nonce++ {
		hash = pow.hash(nonce, &val)
		if val.Cmp(pow.target) == -1 {
			break
		}
	}
	fmt.Printf("# Hash: %x\n\n", hash)
	return nonce, hash[:]
	return nonce, hash
}

func (pow *PoW) Validate() bool {
	var val big.Int
	pow.hash(pow.block.Nonce, &val)
	return val.Cmp(pow.target) == -1
}

func (pow *PoW) hash(nonce int, dst *big.Int) []byte {
	data := pow.prepareData(nonce)
	hash := sha256.Sum256(data)
	dst.SetBytes(hash[:])
	return hash[:]
}

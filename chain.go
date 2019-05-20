package main

type Blockchain struct {
	blocks []*Block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{blocks: []*Block{NewGenesisBlock()}}
}

func (bc *Blockchain) AddBlock(data string) {
	prev := bc.blocks[len(bc.blocks)-1]
	curr := NewBlock(data, prev.Hash)
	bc.blocks = append(bc.blocks, curr)
}

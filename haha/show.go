package haha

import (
	"fmt"
	"haha/blockchain"
	"strconv"
)

func show() {
	bc, err := blockchain.LoadBlockchain()
	if err != nil {
		showError(err)
		return
	}
	defer bc.Close()

	bci := bc.Iterator()
	for {
		block, err := bci.Next()
		if err != nil {
			showError(err)
			break
		}
		if block == nil {
			break
		}

		fmt.Printf("PrevHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewPoW(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}

package haha

import (
	"fmt"
	"github.com/yuanqj/haha-coin/blockchain"
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

		for i := 0; i < 31; i++ {
			fmt.Printf(">")
		}
		fmt.Printf(" Block\n")
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("PrevHash: %x\n", block.PrevBlockHash)
		pow := blockchain.NewPoW(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Printf("Trasactions: \n")
		for _, tx := range block.Transactions {
			fmt.Println(tx.String())
		}
		fmt.Println()
	}
}

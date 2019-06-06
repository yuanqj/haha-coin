package haha

import (
	"fmt"
	"haha/blockchain"
)

func initBlockchain(addr string) {
	bc, err := blockchain.CreateBlockchain(addr)
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.Close()

	fmt.Println("Done!")
}

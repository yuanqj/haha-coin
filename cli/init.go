package cli

import (
	"fmt"
	"haha/blockchain"
)

func (cli *CLI) createBlockchain(address string) {
	bc, err := blockchain.CreateBlockchain(address)
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.Close()

	fmt.Println("Done!")
}

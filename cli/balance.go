package cli

import (
	"fmt"
	"math"
	"haha/blockchain"
)

func (cli *CLI) getBalance(address string) {
	bc, err := blockchain.LoadBlockchain()
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.Close()

	_, tot, err := bc.UTXOs(address, math.MaxInt64)
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.Close()

	fmt.Printf("Balance of '%s': %d\n", address, tot)
}

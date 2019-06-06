package haha

import (
	"fmt"
	"haha/blockchain"
	"haha/wallet"
	"math"
)

func balance(addr string) {
	_, err := wallet.DecodeAddr(addr)
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	bc, err := blockchain.LoadBlockchain()
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.Close()

	_, tot, err := bc.UTXOs(addr, math.MaxInt64)
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.Close()

	fmt.Printf("Balance of '%s': %d\n", addr, tot)
}

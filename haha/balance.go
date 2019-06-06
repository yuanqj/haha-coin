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
		showError(err)
		return
	}
	bc, err := blockchain.LoadBlockchain()
	if err != nil {
		showError(err)
		return
	}
	defer bc.Close()

	_, tot, err := bc.UTXOs(addr, math.MaxInt64)
	if err != nil {
		showError(err)
		return
	}

	fmt.Printf("Balance of '%s': %d\n", addr, tot)
}

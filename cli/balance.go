package cli

import (
	"fmt"
	"math"
	"haha/blockchain"
	"haha/wallet"
)

func (cli *CLI) getBalance(addr string) {
	valid, err := wallet.ValidateAddr(addr)
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	if !valid {
		fmt.Println("************* ERROR: Address is invalid")
		return
	}
	bc, err := blockchain.LoadBlockchain()
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.Close()

	ws, err := wallet.NewWallets()
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}

	_, tot, err := bc.UTXOs(ws.GetWallet(addr), math.MaxInt64)
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.Close()

	fmt.Printf("Balance of '%s': %d\n", addr, tot)
}

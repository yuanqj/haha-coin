package haha

import (
	"fmt"
	"haha/blockchain"
	"haha/wallet"
	"log"
	"math"
)

func wallets() {
	bc, err := blockchain.LoadBlockchain()
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.Close()

	wallets, err := wallet.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addrs := wallets.GetAddrs()

	amts := make([]int, len(addrs))
	for i, addr := range addrs {
		_, tot, err := bc.UTXOs(addr, math.MaxInt64)
		if err != nil {
			fmt.Println("************* Error:")
			fmt.Println(err)
			return
		}
		amts[i] = tot
	}

	for i, addr := range addrs {
		fmt.Printf("%02d: %s, %d\n", i, addr, amts[i])
	}
}

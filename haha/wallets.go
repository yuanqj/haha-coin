package haha

import (
	"fmt"
	"github.com/yuanqj/haha-coin/blockchain"
	"github.com/yuanqj/haha-coin/wallet"
	"math"
)

func wallets() {
	bc, err := blockchain.LoadBlockchain()
	if err != nil {
		showError(err)
		return
	}
	defer bc.Close()

	wallets, err := wallet.NewWallets()
	if err != nil {
		showError(err)
		return
	}
	addrs := wallets.GetAddrs()

	amts := make([]int, len(addrs))
	for i, addr := range addrs {
		_, tot, err := bc.UTXOs(addr, math.MaxInt64)
		if err != nil {
			showError(err)
			return
		}
		amts[i] = tot
	}

	for i, addr := range addrs {
		fmt.Printf("%02d: %s, %d\n", i, addr, amts[i])
	}
}

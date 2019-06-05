package cli

import (
	"log"
	"fmt"
	"haha/wallet"
)

func (cli *CLI) listAddrs() {
	wallets, err := wallet.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddrs()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

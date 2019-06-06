package haha

import (
	"fmt"
	"haha/wallet"
	"log"
)

func wallets() {
	wallets, err := wallet.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addr := wallets.GetAddrs()

	for _, address := range addr {
		fmt.Println(address)
	}
}

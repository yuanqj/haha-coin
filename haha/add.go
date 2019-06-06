package haha

import (
	"fmt"
	"haha/wallet"
)

func add() {
	wallets, _ := wallet.NewWallets()
	address, err := wallets.CreateWallet()
	if err != nil {
		showError(err)
		return
	}
	if err := wallets.Save(); err != nil {
		showError(err)
		return
	}
	fmt.Printf("Your new address: %s\n", address)
}

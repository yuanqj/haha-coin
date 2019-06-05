package cli

import (
	"fmt"
	"haha/wallet"
)

func (cli *CLI) createWallet() {
	wallets, _ := wallet.NewWallets()
	address, err := wallets.CreateWallet()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	if err:= wallets.Save(); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	fmt.Printf("Your new address: %s\n", address)
}

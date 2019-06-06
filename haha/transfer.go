package haha

import (
	"fmt"
	"haha/blockchain"
	"haha/transaction"
	"haha/wallet"
)

func transfer(src, dst string, amt int) {
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
	w := ws.Wallets[src]
	if w == nil {
		fmt.Println("************* Error:")
		fmt.Printf("given address not found in wallets: '%s'\n", src)
		return
	}
	utxos, _, err := bc.UTXOs(src, amt)
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}

	tx, err := transaction.NewUTXOTransaction(w, dst, amt, utxos)
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}

	err = bc.MineBlock([]*transaction.Transaction{tx})
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	fmt.Println("Success!")
}

package cli

import (
	"fmt"
	"haha/blockchain"
	"haha/transaction"
	"haha/wallet"
)

func (cli *CLI) send(from, to string, amt int) {
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
	utxos, _, err := bc.UTXOs(ws.GetWallet(from), amt)
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}

	tx, err := transaction.NewUTXOTransaction(from, to, amt, utxos)
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

package cli

import (
	"fmt"
	"haha/blockchain"
	"haha/transaction"
)

func (cli *CLI) send(from, to string, amount int) {
	bc, err := blockchain.LoadBlockchain()
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.Close()

	tx, err := transaction.NewUTXOTransaction(from, to, amount, bc)
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

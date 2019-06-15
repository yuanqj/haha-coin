package haha

import (
	"fmt"
	"github.com/yuanqj/haha-coin/blockchain"
	"github.com/yuanqj/haha-coin/transaction"
	"github.com/yuanqj/haha-coin/wallet"
)

func transfer(src, dst string, amt int) {
	bc, err := blockchain.LoadBlockchain()
	if err != nil {
		showError(err)
		return
	}
	defer bc.Close()

	ws, err := wallet.NewWallets()
	if err != nil {
		showError(err)
		return
	}
	w := ws.Wallets[src]
	if w == nil {
		showError(fmt.Errorf("given address not found in wallets: '%s'\n", src))
		return
	}
	utxos, _, err := bc.UTXOs(src, amt)
	if err != nil {
		showError(err)
		return
	}

	tx, err := transaction.NewUTXOTransaction(w, dst, amt, utxos)
	if err != nil {
		showError(err)
		return
	}

	err = bc.MineBlock([]*transaction.Transaction{tx})
	if err != nil {
		showError(err)
		return
	}
}

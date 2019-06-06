package haha

import (
	"haha/blockchain"
)

func initBlockchain(addr string) {
	bc, err := blockchain.CreateBlockchain(addr)
	if err != nil {
		showError(err)
		return
	}
	defer bc.Close()
}

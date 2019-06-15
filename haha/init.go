package haha

import (
	"github.com/yuanqj/haha-coin/blockchain"
)

func initBlockchain(addr string) {
	bc, err := blockchain.CreateBlockchain(addr)
	if err != nil {
		showError(err)
		return
	}
	defer bc.Close()
}

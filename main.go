package main

import (
	"fmt"
)

func main() {
	bc, err := NewBlockchain()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"math"
)

// CLI responsible for processing command line arguments
type CLI struct{}

func (cli *CLI) createBlockchain(address string) {
	bc, err := CreateBlockchain(address)
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.db.Close()

	bc.db.Close()
	fmt.Println("Done!")
}

func (cli *CLI) getBalance(address string) {
	bc, err := LoadBlockchain()
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.db.Close()

	_, tot, err := bc.FindUTXOs(address, math.MaxInt64)
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.db.Close()

	fmt.Printf("Balance of '%s': %d\n", address, tot)
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  balance -addr ADDRESS - Get balance of ADDRESS")
	fmt.Println("  init -addr ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  print - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amt AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printChain() {
	bc, err := LoadBlockchain()
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.db.Close()

	bci := bc.Iterator()
	for {
		block, err := bci.Next()
		if err != nil {
			fmt.Println("************* Error:")
			fmt.Println(err)
			break
		}
		if block == nil {
			break
		}

		fmt.Printf("PrevHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewPoW(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}

func (cli *CLI) send(from, to string, amount int) {
	bc, err := LoadBlockchain()
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	defer bc.db.Close()

	tx, err := NewUTXOTransaction(from, to, amount, bc)
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}

	err = bc.MineBlock([]*Transaction{tx})
	if err != nil {
		fmt.Println("************* Error:")
		fmt.Println(err)
		return
	}
	fmt.Println("Success!")
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("init", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("addr", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("addr", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amt", 0, "Amount to send")

	switch os.Args[1] {
	case "balance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "init":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
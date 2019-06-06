package haha

import (
	"github.com/spf13/cobra"
	"fmt"
)

func Run() {

	var cmdAdd = &cobra.Command{
		Use:   "add",
		Short: "add a new wallet",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			add()
		},
	}

	var cmdWallets = &cobra.Command{
		Use:   "wallets",
		Short: "show all wallets",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			wallets()
		},
	}

	var cmdInitAddr string
	var cmdInit = &cobra.Command{
		Use:   "init",
		Short: "initialize the blockchain",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			initBlockchain(cmdInitAddr)
		},
	}
	cmdInit.Flags().StringVarP(&cmdInitAddr, "addr", "", "", "address to receive genesis reward")
	cmdInit.MarkFlagRequired("addr")

	var cmdShow = &cobra.Command{
		Use:   "show",
		Short: "show all blocks in current blockchain",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			show()
		},
	}

	var cmdBalanceAddr string
	var cmdBalance = &cobra.Command{
		Use:   "balance",
		Short: "get spendable balance in given address",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			balance(cmdBalanceAddr)
		},
	}
	cmdBalance.Flags().StringVarP(&cmdBalanceAddr, "addr", "", "", "wallet address to get balance for")
	cmdBalance.MarkFlagRequired("addr")

	var (
		cmdTransferSrc string
		cmdTransferDst string
		cmdTransferAmt int
	)
	var cmdTransfer = &cobra.Command{
		Use:   "transfer",
		Short: "get spendable balance in given address",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			transfer(cmdTransferSrc, cmdTransferDst, cmdTransferAmt)
		},
	}
	cmdTransfer.Flags().StringVarP(&cmdTransferSrc, "src", "", "", "source wallet address")
	cmdTransfer.Flags().StringVarP(&cmdTransferDst, "dst", "", "", "destination wallet address")
	cmdTransfer.Flags().IntVarP(&cmdTransferAmt, "amt", "", 0, "amount of balance to transfer")
	cmdTransfer.MarkFlagRequired("src")
	cmdTransfer.MarkFlagRequired("dst")
	cmdTransfer.MarkFlagRequired("amt")

	var haha = &cobra.Command{Use: "haha"}
	haha.AddCommand(cmdAdd, cmdWallets, cmdInit, cmdShow, cmdBalance, cmdTransfer)
	haha.Execute()
}

func showError(err error) {
	fmt.Println("************* Error:")
	fmt.Println(err)
}

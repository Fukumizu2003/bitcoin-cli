/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bitcoin-cli/internal/config"
	"bitcoin-cli/internal/util"
	"fmt"

	"github.com/spf13/cobra"
)

var showBalance bool
var showAddress bool
var all bool

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if showBalance {
			err := util.RefreshUtxos()
			if err != nil {
				return err
			}
		}
		if all {
			if showAddress {
				accounts := util.LoadAccounts()
				accounts = append(accounts, util.LoadDestinations()...)
				util.ShowAllAddress(accounts)
			} else if showBalance {
				accounts := util.LoadAccounts()
				utxos := util.LoadUtxos()
				util.ShowAllBalance(accounts, utxos)
			} else {
				return fmt.Errorf("フラグを指定してください。\nアドレス表示： -a\n残高表示： -b")
			}
		} else {
			if showAddress {
				ac := config.GetMainAccount()
				fmt.Println(ac.Address)
				return util.ShowQRCode(ac.Address)
			} else if showBalance {
				utxos := util.LoadUtxos()
				ac := config.GetMainAccount()
				balances := util.GetBalanceBook(utxos)
				fmt.Println(util.SatsToBtc(util.IntToStr(balances[ac.Address])) + " BTC")
			} else {
				return fmt.Errorf("フラグを指定してください。\nアドレス表示： -a\n残高表示： -b")
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVarP(&showBalance, "balance", "b", false, "アカウントの残高を表示")
	showCmd.Flags().BoolVarP(&showAddress, "address", "a", false, "アカウントのアドレスを表示")
	showCmd.Flags().BoolVar(&all, "all", false, "すべてのアカウントの情報を見たい場合にセット")
}

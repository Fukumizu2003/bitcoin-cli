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

var show_balance bool
var show_address bool
var all bool

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if show_balance {
			err := util.Refresh_utxos()
			if err != nil {
				return err
			}
		}
		if all {
			if show_address {
				accounts := util.Load_accounts()
				accounts = append(accounts, util.Load_destinations()...)
				util.Show_all_address(accounts)
			} else if show_balance {
				accounts := util.Load_accounts()
				utxos := util.Load_utxos()
				util.Show_all_balance(accounts, utxos)
			} else {
				return fmt.Errorf("フラグを指定してください。\nアドレス表示： -a\n残高表示： -b")
			}
		} else {
			if show_address {
				ac := config.Get_main_account()
				fmt.Println(ac.Address)
				return util.ShowQRCode(ac.Address)
			} else if show_balance {
				utxos := util.Load_utxos()
				ac := config.Get_main_account()
				balances := util.Get_balance_book(utxos)
				fmt.Println(util.Sats_to_btc(util.Int_to_str(balances[ac.Address])) + " BTC")
			} else {
				return fmt.Errorf("フラグを指定してください。\nアドレス表示： -a\n残高表示： -b")
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVarP(&show_balance, "balance", "b", false, "アカウントの残高を表示")
	showCmd.Flags().BoolVarP(&show_address, "address", "a", false, "アカウントのアドレスを表示")
	showCmd.Flags().BoolVar(&all, "all", false, "すべてのアカウントの情報を見たい場合にセット")
}

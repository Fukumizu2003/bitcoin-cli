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

var acname string

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "メインアカウントをセットします。",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if acname == "" {
			return fmt.Errorf("アカウント名を-nで指定してください。")
		}
		address, err := util.GetAddressFromName(acname)
		if err != nil {
			return fmt.Errorf("このアカウント名は存在しません。")
		}
		config.ChangeMainAccount(address)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&acname, "name", "n", "", "アカウント名によりメインアカウントを設定します。")
}

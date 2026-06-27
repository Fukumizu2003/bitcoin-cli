/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bitcoin-cli/internal/config"
	"fmt"

	"github.com/spf13/cobra"
)

var acname string

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "メインアカウントをセットします。",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if acname == "" {
			fmt.Println("アカウント名を-nで指定してください。")
			return
		}

		st, er := config.SetAccount(acname)
		if er != nil {
			fmt.Println(er)
			return
		}
		config.SaveConfig(*st)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&acname, "name", "n", "", "アカウント名によりメインアカウントを設定します。")
}

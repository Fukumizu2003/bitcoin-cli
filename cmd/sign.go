/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bitcoin-cli/internal/config"
	"bitcoin-cli/internal/util"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
)

var password string

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if password == "" {
			return fmt.Errorf("パスワードを-pにより入力してください。")
		}
		tx := util.LoadTx()
		selfAc := config.GetAccount()
		privkeyCrB64 := selfAc.Key
		privkeyCr := util.B64Decode(privkeyCrB64)
		privkey := util.AesDecrypt(privkeyCr, []byte(password))
		tx = util.GetSignedTx(tx, privkey)
		rawtx := util.TxToBytes(tx)
		fmt.Println("Raw tx hex:")
		fmt.Println(hex.EncodeToString(rawtx))
		fmt.Println("\nTxID (Big endian): " + hex.EncodeToString(util.TxToTxid(tx)))
		util.SaveTx(tx)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringVarP(&password, "password", "p", "", "Set password")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

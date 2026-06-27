/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bitcoin-cli/internal/util"
	"encoding/hex"
	"fmt"
	"slices"

	"github.com/spf13/cobra"
)

var otherAPI bool

// broadcastCmd represents the broadcast command
var broadcastCmd = &cobra.Command{
	Use:   "broadcast",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		tx := util.LoadTx()
		txidb := util.TxToTxid(tx)
		slices.Reverse(txidb)
		txid := hex.EncodeToString(txidb)
		rawTx := util.TxToBytes(tx)
		var msg []byte
		var err error
		if !otherAPI {
			msg, err = util.Broadcast(rawTx)
		} else {
			msg, err = util.Broadcast2(rawTx)
		}
		if err != nil {
			return err
		}
		if string(msg) == txid {
			fmt.Println("SUCCEED: " + string(msg))
			return nil
		}
		fmt.Println("API response: " + string(msg))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(broadcastCmd)
	broadcastCmd.Flags().BoolVar(&otherAPI, "another", false, "")
}

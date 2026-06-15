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
		rawTx := util.TxToBytes(tx)
		txid := util.TxToTxid(tx)
		slices.Reverse(txid)
		fmt.Println(hex.EncodeToString(rawTx))
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
		msgId, _ := hex.DecodeString(string(msg))
		if slices.Equal(txid, msgId) {
			fmt.Println("SUCCEED: " + hex.EncodeToString(txid))
			return nil
		}
		fmt.Println("Something went wrong.")
		fmt.Println(string(msg))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(broadcastCmd)
	broadcastCmd.Flags().BoolVar(&otherAPI, "another", false, "")
}

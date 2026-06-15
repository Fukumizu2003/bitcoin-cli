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
		tx := util.Load_tx()
		raw_tx := util.Tx_to_bytes(tx)
		txid := util.Tx_to_txid(tx)
		slices.Reverse(txid)
		fmt.Println(hex.EncodeToString(raw_tx))
		var msg []byte
		var err error
		if !otherAPI {
			msg, err = util.Broadcast(raw_tx)
		} else {
			msg, err = util.Broadcast2(raw_tx)
		}
		if err != nil {
			return err
		}
		msg_id, _ := hex.DecodeString(string(msg))
		if slices.Equal(txid, msg_id) {
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

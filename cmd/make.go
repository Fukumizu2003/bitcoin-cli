/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bitcoin-cli/internal/config"
	"bitcoin-cli/internal/util"
	"fmt"
	"strconv"
	"strings"

	"encoding/binary"
	"encoding/hex"

	"github.com/spf13/cobra"
)

var version_str string
var amount_str string
var destination string
var fee_str string

// makeCmd represents the make command
var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "トランザクションを設定します。署名は行わないのでオンライン状態で大丈夫です。",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		tx := util.New_tx()
		ac := config.Get_main_account()
		tx.Senders = append(tx.Senders, ac.Address)

		if destination == "" {
			return fmt.Errorf("送金先を指定してください。")
		}
		if amount_str == "" {
			return fmt.Errorf("送金額を指定してください。")
		}

		dests := strings.Split(destination, " ")
		destinations := []string{}
		for _, dest := range dests {
			addr, err := util.Get_address_from_name(dest)
			if err != nil {
				destinations = append(destinations, dest)
			} else {
				destinations = append(destinations, addr)
			}
		}
		amounts_btc := strings.Split(amount_str, " ")
		amounts_int := []int{}
		for _, btc := range amounts_btc {
			sats := util.Btc_to_sats(btc)
			amounts_int = append(amounts_int, sats)
		}

		all_utxos := util.Load_utxos()
		utxos := all_utxos[ac.Address]

		version_int := util.Str_to_int(version_str)
		amount_int := 0
		for _, am := range amounts_int {
			amount_int += am
		}
		fee_int, _ := strconv.Atoi(fee_str)
		utxos_consume, err := util.Necessary_inputs(utxos, amount_int+fee_int)

		if err != nil {
			return err
		}

		ver, _ := util.Int_to_bytes(version_int, 4)
		tx.Version = ver
		tx.Marker = []byte{0x00}
		tx.Flag = []byte{0x01}
		tx.Inputcount, _ = util.Int_to_compactsize(len(utxos_consume))
		for _, ut := range utxos_consume {
			tx.Inputs = append(tx.Inputs, util.Utxo_to_input(ut))
			witness := util.New_witness()
			tx.Witness = append(tx.Witness, witness)
		}
		total_input := util.Total_of_utxos(utxos_consume)
		margin := total_input - amount_int - fee_int
		if margin <= 546 {
			tx.Outputcount, _ = util.Int_to_compactsize(len(destinations))
			for i, dest := range destinations {
				output := util.New_output()
				spk, _ := util.Scriptpubkey(dest)
				output.Amount, _ = util.Int_to_bytes(amounts_int[i], 8)
				output.Scriptpubkeysize, _ = util.Int_to_compactsize(len(spk))
				output.Scriptpubkey = spk
				tx.Outputs = append(tx.Outputs, output)
			}
		} else {
			tx.Outputcount, _ = util.Int_to_compactsize(len(destinations) + 1)
			for i, dest := range destinations {
				output := util.New_output()
				spk, _ := util.Scriptpubkey(dest)
				output.Amount, _ = util.Int_to_bytes(amounts_int[i], 8)
				output.Scriptpubkeysize, _ = util.Int_to_compactsize(len(spk))
				output.Scriptpubkey = spk
				tx.Outputs = append(tx.Outputs, output)
			}
			output := util.New_output()
			spk, _ := util.Scriptpubkey(ac.Address)
			output.Amount, _ = util.Int_to_bytes(margin, 8)
			output.Scriptpubkeysize, _ = util.Int_to_compactsize(len(spk))
			output.Scriptpubkey = spk
			tx.Outputs = append(tx.Outputs, output)
		}
		fee_confirm := 0
		for _, inp := range tx.Inputs {
			fee_confirm += util.Get_utxos_value(all_utxos, hex.EncodeToString(inp.Txid), hex.EncodeToString(inp.Vout))
		}
		for _, oup := range tx.Outputs {
			fee_confirm -= int(binary.LittleEndian.Uint64(oup.Amount))
		}
		fmt.Println("Fee: " + strconv.Itoa(fee_confirm) + " sats")
		tx.Locktime = []byte{0x00, 0x00, 0x00, 0x00}
		util.Save_tx(tx)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(makeCmd)

	makeCmd.Flags().StringVarP(&version_str, "version", "v", "2", "")
	makeCmd.Flags().StringVarP(&amount_str, "amount", "a", "", "着金額を指定。複数指定する場合はアドレスに対応した順序で半角スペース区切りで並べる。")
	makeCmd.Flags().StringVarP(&destination, "destination", "d", "", "送信先アドレスを指定。複数指定する場合は半角スペース区切りで並べる。")
	makeCmd.Flags().StringVarP(&fee_str, "fee", "f", "500", "Transaction fee by sats")
}

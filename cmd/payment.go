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

var versionStr string
var amountStr string
var destination string
var feeStr string

// paymentCmd represents the payment command
var paymentCmd = &cobra.Command{
	Use:   "payment",
	Short: "トランザクションを設定します。署名は行わないのでオンライン状態で大丈夫です。",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		tx := util.NewTx()
		ac := config.GetMainAccount()
		tx.Senders = append(tx.Senders, ac.Address)

		if destination == "" {
			return fmt.Errorf("送金先を指定してください。")
		}
		if amountStr == "" {
			return fmt.Errorf("送金額を指定してください。")
		}

		dests := strings.Split(destination, " ")
		destinations := []string{}
		for _, dest := range dests {
			addr, err := util.GetAddressFromName(dest)
			if err != nil {
				destinations = append(destinations, dest)
			} else {
				destinations = append(destinations, addr)
			}
		}
		amountsBtc := strings.Split(amountStr, " ")
		amountsInt := []int{}
		for _, btc := range amountsBtc {
			sats := util.BtcToSats(btc)
			amountsInt = append(amountsInt, sats)
		}

		allUtxos := util.LoadUtxos()
		utxos := allUtxos[ac.Address]

		versionInt := util.StrToInt(versionStr)
		amountInt := 0
		for _, am := range amountsInt {
			amountInt += am
		}
		feeInt, _ := strconv.Atoi(feeStr)
		utxosConsume, err := util.NecessaryInputs(utxos, amountInt+feeInt)

		if err != nil {
			return err
		}

		ver, _ := util.IntToBytes(versionInt, 4)
		tx.Version = ver
		tx.Marker = []byte{0x00}
		tx.Flag = []byte{0x01}
		tx.Inputcount, _ = util.IntToCompactsize(len(utxosConsume))
		for _, ut := range utxosConsume {
			tx.Inputs = append(tx.Inputs, util.UtxoToInput(ut))
			witness := util.NewWitness()
			tx.Witness = append(tx.Witness, witness)
		}
		totalInput := util.TotalOfUtxos(utxosConsume)
		margin := totalInput - amountInt - feeInt
		if margin <= 546 {
			tx.Outputcount, _ = util.IntToCompactsize(len(destinations))
			for i, dest := range destinations {
				output := util.NewOutput()
				spk, _ := util.Scriptpubkey(dest)
				output.Amount, _ = util.IntToBytes(amountsInt[i], 8)
				output.Scriptpubkeysize, _ = util.IntToCompactsize(len(spk))
				output.Scriptpubkey = spk
				tx.Outputs = append(tx.Outputs, output)
			}
		} else {
			tx.Outputcount, _ = util.IntToCompactsize(len(destinations) + 1)
			for i, dest := range destinations {
				output := util.NewOutput()
				spk, _ := util.Scriptpubkey(dest)
				output.Amount, _ = util.IntToBytes(amountsInt[i], 8)
				output.Scriptpubkeysize, _ = util.IntToCompactsize(len(spk))
				output.Scriptpubkey = spk
				tx.Outputs = append(tx.Outputs, output)
			}
			output := util.NewOutput()
			spk, _ := util.Scriptpubkey(ac.Address)
			output.Amount, _ = util.IntToBytes(margin, 8)
			output.Scriptpubkeysize, _ = util.IntToCompactsize(len(spk))
			output.Scriptpubkey = spk
			tx.Outputs = append(tx.Outputs, output)
		}
		feeConfirm := 0
		for _, inp := range tx.Inputs {
			feeConfirm += util.GetUtxosValue(allUtxos, hex.EncodeToString(inp.Txid), hex.EncodeToString(inp.Vout))
		}
		for _, oup := range tx.Outputs {
			feeConfirm -= int(binary.LittleEndian.Uint64(oup.Amount))
		}
		fmt.Println("Fee: " + strconv.Itoa(feeConfirm) + " sats")
		tx.Locktime = []byte{0x00, 0x00, 0x00, 0x00}
		util.SaveTx(tx)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(paymentCmd)

	paymentCmd.Flags().StringVarP(&versionStr, "version", "v", "2", "")
	paymentCmd.Flags().StringVarP(&amountStr, "amount", "a", "", "着金額を指定。複数指定する場合はアドレスに対応した順序で半角スペース区切りで並べる。")
	paymentCmd.Flags().StringVarP(&destination, "destination", "d", "", "送信先アドレスを指定。複数指定する場合は半角スペース区切りで並べる。")
	paymentCmd.Flags().StringVarP(&feeStr, "fee", "f", "500", "Transaction fee by sats")
}

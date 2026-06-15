package util

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

type Input struct {
	Txid          []byte
	Vout          []byte
	Scriptsigsize []byte
	Scriptsig     []byte
	Sequence      []byte
}

type Output struct {
	Amount           []byte
	Scriptpubkeysize []byte
	Scriptpubkey     []byte
}

type Witness struct {
	Stackitems []byte
	Items      [][]byte
}

type Tx struct {
	Senders     []string
	Version     []byte
	Marker      []byte
	Flag        []byte
	Inputcount  []byte
	Inputs      []Input
	Outputcount []byte
	Outputs     []Output
	Witness     []Witness
	Locktime    []byte
}

type Balance struct {
	Vout  []byte
	Value []byte
}

func New_tx() Tx {
	tx := Tx{
		Senders:     make([]string, 0, 2),
		Version:     make([]byte, 0, 4),
		Marker:      make([]byte, 0, 1),
		Flag:        make([]byte, 0, 1),
		Inputcount:  make([]byte, 0, 1),
		Inputs:      make([]Input, 0, 4),
		Outputcount: make([]byte, 0, 1),
		Outputs:     make([]Output, 0, 4),
		Witness:     make([]Witness, 0, 4),
		Locktime:    make([]byte, 0, 4),
	}
	return tx
}

func New_input() Input {
	inp := Input{
		Txid:          make([]byte, 0, 32),
		Vout:          make([]byte, 0, 4),
		Scriptsigsize: make([]byte, 0, 1),
		Scriptsig:     make([]byte, 0, 100),
		Sequence:      make([]byte, 0, 4),
	}
	return inp
}

func New_output() Output {
	oup := Output{
		Amount:           make([]byte, 0, 8),
		Scriptpubkeysize: make([]byte, 0, 1),
		Scriptpubkey:     make([]byte, 0, 25),
	}
	return oup
}

func New_witness() Witness {
	items := make([][]byte, 0, 4)
	wit := Witness{
		Stackitems: make([]byte, 0, 1),
		Items:      items,
	}
	return wit
}

func Save_tx(tx Tx) {
	Mkdir_or_nothing("temp")
	data, err := json.MarshalIndent(tx, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile(Relative_to_absolute("temp", "transaction.json"), data, 0644)
}

func Load_tx() Tx {
	data, _ := os.ReadFile(Relative_to_absolute("temp", "transaction.json"))
	var tx Tx
	json.Unmarshal(data, &tx)
	return tx
}

/*

func Sort_input(inputs []map[string]string) []map[string]string {
	res := []map[string]string{}
	values := []int{}
	for _, input := range inputs {
		values = append(values, Str_to_int(input["value"]))
	}
	sorted_indexes := Get_sorted_indexes(values)
	for _, idx := range sorted_indexes {
		res = append(res, inputs[idx])
	}
	return res
}

*/

func Utxo_to_input(utxo map[string]string) Input {
	input := New_input()
	input.Txid, _ = hex.DecodeString(utxo["txid"])
	input.Vout, _ = hex.DecodeString(utxo["vout"])
	input.Scriptsigsize = []byte{0x00}
	input.Sequence = []byte{0xff, 0xff, 0xff, 0xff}
	return input
}

func Total_of_utxos(utxos []map[string]string) int {
	total := 0
	for _, ut := range utxos {
		total += Str_to_int(ut["value"])
	}
	return total
}

func Tx_to_bytes(tx Tx) []byte {
	msg := []byte{}
	msg = append(msg, tx.Version...)
	msg = append(msg, tx.Marker...)
	msg = append(msg, tx.Flag...)
	msg = append(msg, tx.Inputcount...)
	for _, input := range tx.Inputs {
		msg = append(msg, input.Txid...)
		msg = append(msg, input.Vout...)
		msg = append(msg, input.Scriptsigsize...)
		msg = append(msg, input.Scriptsig...)
		msg = append(msg, input.Sequence...)
	}
	msg = append(msg, tx.Outputcount...)
	for _, output := range tx.Outputs {
		msg = append(msg, output.Amount...)
		msg = append(msg, output.Scriptpubkeysize...)
		msg = append(msg, output.Scriptpubkey...)
	}
	for _, witness := range tx.Witness {
		msg = append(msg, witness.Stackitems...)
		for _, item := range witness.Items {
			msg = append(msg, item...)
		}
	}
	msg = append(msg, tx.Locktime...)
	return msg
}

func Tx_to_txid(tx Tx) []byte {
	msg := []byte{}
	msg = append(msg, tx.Version...)
	msg = append(msg, tx.Inputcount...)
	for _, input := range tx.Inputs {
		msg = append(msg, input.Txid...)
		msg = append(msg, input.Vout...)
		msg = append(msg, input.Scriptsigsize...)
		msg = append(msg, input.Scriptsig...)
		msg = append(msg, input.Sequence...)
	}
	msg = append(msg, tx.Outputcount...)
	for _, output := range tx.Outputs {
		msg = append(msg, output.Amount...)
		msg = append(msg, output.Scriptpubkeysize...)
		msg = append(msg, output.Scriptpubkey...)
	}
	for _, witness := range tx.Witness {
		msg = append(msg, witness.Stackitems...)
		for _, item := range witness.Items {
			msg = append(msg, item...)
		}
	}
	msg = append(msg, tx.Locktime...)
	return Hash256(msg)
}

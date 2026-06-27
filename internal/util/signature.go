package util

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/bech32"
)

func Scriptpubkey(address string) ([]byte, error) {
	adtype, _ := AddressType(address)
	if adtype == "p2pkh" {
		decoded, _, err := base58.CheckDecode(address)
		if err != nil {
			return nil, err
		}
		head := []byte{0x76, 0xa9, 0x14}
		tail := []byte{0x88, 0xac}
		return append(append(head, decoded...), tail...), nil
	} else if adtype == "p2wpkh" {
		head := []byte{0x00, 0x14}
		_, decoded, err := bech32.Decode(address)
		if err != nil {
			fmt.Println(err)
		}
		converted, _ := bech32.ConvertBits(decoded[1:], 5, 8, false)
		return append(head, converted...), nil
	} else if adtype == "p2sh" {
		dec, _, err := base58.CheckDecode(address)
		if err != nil {
			return nil, err
		}
		head := []byte{0xa9, 0x14}
		tail := []byte{0x87}
		return append(append(head, dec...), tail...), nil
	} else {
		return nil, errors.New("未対応のアドレス形式")
	}
}

func MakeSignatureFromBytes(priv []byte, msg []byte) []byte {
	privkey, _ := BytesToKeypair(priv)
	hash := Hash256(msg)
	signature := ecdsa.Sign(privkey, hash)
	return append(signature.Serialize(), 0x01)
}

func getPubkey(priv []byte) []byte {
	_, pubkey := BytesToKeypair(priv)
	return pubkey.SerializeCompressed()
}

func MakeSignatureForSegwitInput(tx Tx, priv []byte, inputindex int) ([]byte, error) {
	pub := getPubkey(priv)

	input := tx.Inputs[inputindex]
	txid := input.Txid
	vout := input.Vout
	sequence := input.Sequence

	utxoData := LoadUtxos()
	balanceInt := GetUtxosValue(utxoData, hex.EncodeToString(txid), hex.EncodeToString(vout))
	if balanceInt == 0 {
		return nil, errors.New("Failed to get utxo data.")
	}
	balance, _ := IntToBytes(balanceInt, 8)
	msg := []byte{}
	msg = append(msg, tx.Version...)
	prehashprevouts := []byte{}
	prehashsequence := []byte{}
	for _, input := range tx.Inputs {
		prehashprevouts = append(append(prehashprevouts, input.Txid...), input.Vout...)
		prehashsequence = append(prehashsequence, input.Sequence...)
	}
	msg = append(msg, Hash256(prehashprevouts)...)
	msg = append(msg, Hash256(prehashsequence)...)
	msg = append(append(msg, txid...), vout...)
	sccode := append(
		append(
			[]byte{0x19, 0x76, 0xa9, 0x14},
			Hash160(pub)...,
		),
		[]byte{0x88, 0xac}...,
	)
	msg = append(msg, sccode...)
	msg = append(msg, balance...)
	msg = append(msg, sequence...)

	prehashoutputs := []byte{}
	for _, output := range tx.Outputs {
		prehashoutputs = append(prehashoutputs, output.Amount...)
		prehashoutputs = append(prehashoutputs, output.Scriptpubkeysize...)
		prehashoutputs = append(prehashoutputs, output.Scriptpubkey...)
	}
	hashoutputs := Hash256(prehashoutputs)
	msg = append(msg, hashoutputs...)

	msg = append(msg, tx.Locktime...)
	msg = append(msg, []byte{0x01, 0x00, 0x00, 0x00}...)
	signature := MakeSignatureFromBytes(priv, msg)
	return signature, nil
}

func GetSignedTx(tx Tx, priv []byte) Tx {
	signatures := [][]byte{}
	pubkey := getPubkey(priv)
	for i := 0; i < int(tx.Inputcount[0]); i++ {
		sign, _ := MakeSignatureForSegwitInput(tx, priv, i)
		signatures = append(signatures, sign)
	}
	for i, signature := range signatures {
		tx.Witness[i].Stackitems = []byte{0x02}
		item0 := append([]byte{byte(len(signature))}, signature...)
		item1 := append([]byte{byte(len(pubkey))}, pubkey...)
		tx.Witness[i].Items = append(tx.Witness[i].Items, item0)
		tx.Witness[i].Items = append(tx.Witness[i].Items, item1)
	}
	return tx
}

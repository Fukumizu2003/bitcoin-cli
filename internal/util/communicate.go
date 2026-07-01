package util

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"
)

func getUtxos(address string) ([]map[string]string, error) {
	var result []map[string]string
	url := "https://blockstream.info/api/address/" + address + "/utxo"
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data []map[string]interface{}
	json.NewDecoder(res.Body).Decode(&data)
	for _, utxo := range data {
		txid := utxo["txid"].(string)
		txidBytes, _ := hex.DecodeString(txid)
		slices.Reverse(txidBytes)

		vout := uint32(utxo["vout"].(float64))
		value := int(utxo["value"].(float64))

		voutBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(voutBytes, vout)

		d := make(map[string]string)
		d["txid"] = hex.EncodeToString(txidBytes)
		d["vout"] = hex.EncodeToString(voutBytes)
		d["value"] = strconv.Itoa(value)
		result = append(result, d)
	}
	return result, nil
}

func RefreshUtxos() error {
	MkdirOrNothing("ref")
	data := LoadAccounts()

	toSaveMap := make(map[string][]map[string]string)
	addresses := []string{}
	for _, ac := range data {
		addresses = append(addresses, ac[1])
	}
	for _, address := range addresses {
		utxos, err := getUtxos(address)
		if err != nil {
			return err
		}
		toSaveMap[address] = utxos
		time.Sleep(500 * time.Millisecond)
	}
	toSave, _ := json.MarshalIndent(toSaveMap, "", "    ")
	os.WriteFile(RelativeToAbsolute("ref", "BTC_utxos.json"), toSave, 0644)
	return nil
}

func Broadcast(raw []byte) ([]byte, error) {
	url := "https://blockstream.info/api/tx"
	rawHex := hex.EncodeToString(raw)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(rawHex))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	SaveResp(body)
	return body, err
}

func Broadcast2(raw []byte) ([]byte, error) {
	url := "https://mempool.space/api/tx"
	rawHex := hex.EncodeToString(raw)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(rawHex))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	SaveResp(body)
	return body, err
}

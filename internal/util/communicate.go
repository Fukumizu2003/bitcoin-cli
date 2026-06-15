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

func get_utxos(address string) ([]map[string]string, error) {
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
		txid_bytes, _ := hex.DecodeString(txid)
		slices.Reverse(txid_bytes)

		vout := uint32(utxo["vout"].(float64))
		value := int(utxo["value"].(float64))

		vout_bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(vout_bytes, vout)

		d := make(map[string]string)
		d["txid"] = hex.EncodeToString(txid_bytes)
		d["vout"] = hex.EncodeToString(vout_bytes)
		d["value"] = strconv.Itoa(value)
		result = append(result, d)
	}
	return result, nil
}

func Refresh_utxos() error {
	Mkdir_or_nothing("ref")
	data := Load_accounts()

	to_save_map := make(map[string][]map[string]string)
	addresses := []string{}
	for _, ac := range data {
		addresses = append(addresses, ac[1])
	}
	for _, address := range addresses {
		utxos, err := get_utxos(address)
		if err != nil {
			return err
		}
		to_save_map[address] = utxos
		time.Sleep(500 * time.Millisecond)
	}
	to_save, _ := json.MarshalIndent(to_save_map, "", "    ")
	os.WriteFile(Relative_to_absolute("ref", "utxos.json"), to_save, 0644)
	return nil
}

func Broadcast(raw []byte) ([]byte, error) {
	url := "https://blockstream.info/api/tx"
	raw_hex := hex.EncodeToString(raw)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(raw_hex))
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
	return body, err
}

func Broadcast2(raw []byte) ([]byte, error) {
	url := "https://mempool.space/api/tx"
	raw_hex := hex.EncodeToString(raw)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(raw_hex))
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
	return body, err
}

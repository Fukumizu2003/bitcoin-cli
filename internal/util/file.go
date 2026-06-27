package util

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

func MkdirOrNothing(dir string) {
	os.MkdirAll(dir, 0755)
}

func LoadAccounts() [][]string {
	MkdirOrNothing("ref")
	f, _ := os.Open(RelativeToAbsolute("ref", "keypair.csv"))
	defer f.Close()
	data, _ := csv.NewReader(f).ReadAll()
	return data
}

func LoadDestinations() [][]string {
	MkdirOrNothing("ref")
	f, _ := os.Open(RelativeToAbsolute("ref", "destinations.csv"))
	defer f.Close()
	data, _ := csv.NewReader(f).ReadAll()
	return data
}

func GetAddressFromName(name string) (string, error) {
	accounts := LoadAccounts()
	accounts = append(accounts, LoadDestinations()...)
	addr := ""
	flag := false
	for _, ac := range accounts {
		if ac[0] == name {
			addr = ac[1]
			flag = true
			break
		}
	}
	if !flag {
		return "", errors.New("指定のアカウント名は存在しません。")
	}
	return addr, nil
}

func LoadUtxos() map[string][]map[string]string {
	MkdirOrNothing("ref")
	data, _ := os.ReadFile(RelativeToAbsolute("ref", "utxos.json"))
	var bn map[string][]map[string]string
	json.Unmarshal(data, &bn)
	return bn
}

func GetAddressUtxos(utxos map[string][]map[string]string, address string) []map[string]string {
	return utxos[address]
}

func GetUtxosValue(utxos map[string][]map[string]string, txid string, vout string) int {
	value := 0
	for _, uos := range utxos {
		for _, uo := range uos {
			if txid == uo["txid"] && vout == uo["vout"] {
				value, _ = strconv.Atoi(uo["value"])
			}
		}
	}
	return value
}

func GetBalanceBook(utxos map[string][]map[string]string) map[string]int {
	res := make(map[string]int)
	for address, utxos := range utxos {
		balance := 0
		for _, val := range utxos {
			v, _ := strconv.Atoi(val["value"])
			balance += v
		}
		res[address] = balance
	}
	return res
}

/*
func LoadUtxosHex() map[string][]map[string]string {
	bn := LoadUtxos()
	res := make(map[string][]map[string]string)
	for address, _ := range bn {
		res[address] = []map[string]string{}
	}

	for address, info := range bn {
		for _, val := range info {
			slices.Reverse(val["txid"])
			txid := hex.EncodeToString(val["txid"])
			vout := strconv.Itoa(int(binary.LittleEndian.Uint32(val["vout"])))
			value := strconv.Itoa(int(binary.LittleEndian.Uint64(val["value"])))

			utxo := make(map[string]string)
			utxo["txid"] = txid
			utxo["vout"] = vout
			utxo["value"] = value

			res[address] = append(res[address], utxo)
		}
	}

	return res
}
*/

func SaveKeypair(acname string, address string, priv []byte) {
	MkdirOrNothing("ref")
	privB64 := B64Encode(priv)
	f, _ := os.OpenFile(RelativeToAbsolute("ref", "keypair.csv"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	row := []byte{}
	row = append(row, []byte(acname)...)
	row = append(row, byte(','))
	row = append(row, []byte(address)...)
	row = append(row, byte(','))
	row = append(row, []byte(privB64)...)
	row = append(row, []byte("\n")...)
	f.Write(row)
}

func SaveAddress(acname string, address string) {
	MkdirOrNothing("ref")
	f, _ := os.OpenFile(RelativeToAbsolute("ref", "destinations.csv"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	row := []byte{}
	row = append(row, []byte(acname)...)
	row = append(row, byte(','))
	row = append(row, []byte(address)...)
	row = append(row, byte('\n'))
	f.Write(row)
}

func SaveResp(data []byte) {
	MkdirOrNothing("temp")
	os.WriteFile(RelativeToAbsolute("temp"), data, 0644)
}

func CheckName(acs [][]string, dss [][]string, name string) bool {
	for _, ac := range acs {
		if ac[0] == name {
			return false
		}
	}
	for _, ds := range dss {
		if ds[0] == name {
			return false
		}
	}
	return true
}

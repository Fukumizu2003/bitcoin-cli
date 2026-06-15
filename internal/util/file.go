package util

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

func Mkdir_or_nothing(dir string) {
	os.MkdirAll(dir, 0755)
}

func Load_accounts() [][]string {
	Mkdir_or_nothing("ref")
	f, _ := os.Open(Relative_to_absolute("ref", "keypair.csv"))
	defer f.Close()
	data, _ := csv.NewReader(f).ReadAll()
	return data
}

func Load_destinations() [][]string {
	Mkdir_or_nothing("ref")
	f, _ := os.Open(Relative_to_absolute("ref", "destinations.csv"))
	defer f.Close()
	data, _ := csv.NewReader(f).ReadAll()
	return data
}

func Get_address_from_name(name string) (string, error) {
	accounts := Load_accounts()
	accounts = append(accounts, Load_destinations()...)
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

func Load_utxos() map[string][]map[string]string {
	Mkdir_or_nothing("ref")
	data, _ := os.ReadFile(Relative_to_absolute("ref", "utxos.json"))
	var bn map[string][]map[string]string
	json.Unmarshal(data, &bn)
	return bn
}

func Get_address_utxos(utxos map[string][]map[string]string, address string) []map[string]string {
	return utxos[address]
}

func Get_utxos_value(utxos map[string][]map[string]string, txid string, vout string) int {
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

func Get_balance_book(utxos map[string][]map[string]string) map[string]int {
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
func Load_utxos_hex() map[string][]map[string]string {
	bn := Load_utxos()
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

func Save_keypair(acname string, address string, priv []byte) {
	Mkdir_or_nothing("ref")
	priv_b64 := B64_encode(priv)
	f, _ := os.OpenFile(Relative_to_absolute("ref", "keypair.csv"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	row := []byte{}
	row = append(row, []byte(acname)...)
	row = append(row, byte(','))
	row = append(row, []byte(address)...)
	row = append(row, byte(','))
	row = append(row, []byte(priv_b64)...)
	row = append(row, []byte("\n")...)
	f.Write(row)
}

func Save_address(acname string, address string) {
	Mkdir_or_nothing("ref")
	f, _ := os.OpenFile(Relative_to_absolute("ref", "destinations.csv"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	row := []byte{}
	row = append(row, []byte(acname)...)
	row = append(row, byte(','))
	row = append(row, []byte(address)...)
	row = append(row, byte('\n'))
	f.Write(row)
}

func Check_name(acs [][]string, dss [][]string, name string) bool {
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

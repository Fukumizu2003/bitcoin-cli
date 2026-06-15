package util

import (
	"fmt"
	"strconv"
)

func Show_all_address(accounts [][]string) {
	fmt.Println("")
	names := []string{}
	addresses := []string{}
	for _, ac := range accounts {
		names = append(names, ac[0])
		addresses = append(addresses, ac[1])
	}
	for i, name := range names {
		fmt.Println(name + ": " + addresses[i])
	}
}

func Show_all_balance(accounts [][]string, utxos map[string][]map[string]string) {
	fmt.Println("")
	ad_name_map := make(map[string]string)
	balancebook := Get_balance_book(utxos)
	for _, ac := range accounts {
		ad_name_map[ac[1]] = ac[0]
	}
	for address, balance := range balancebook {
		btc := Sats_to_btc(strconv.Itoa(balance))
		fmt.Println(ad_name_map[address] + ": " + btc + " BTC")
	}
}

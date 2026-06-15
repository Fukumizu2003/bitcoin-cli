package util

import (
	"fmt"
	"strconv"
)

func ShowAllAddress(accounts [][]string) {
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

func ShowAllBalance(accounts [][]string, utxos map[string][]map[string]string) {
	fmt.Println("")
	adNameMap := make(map[string]string)
	balancebook := GetBalanceBook(utxos)
	for _, ac := range accounts {
		adNameMap[ac[1]] = ac[0]
	}
	for address, balance := range balancebook {
		btc := SatsToBtc(strconv.Itoa(balance))
		fmt.Println(adNameMap[address] + ": " + btc + " BTC")
	}
}

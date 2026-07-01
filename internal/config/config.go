package config

import (
	"bitcoin-cli/internal/util"
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type State struct {
	Name    string
	Address string
	Key     string
}

func SetAccount(name string) (*State, error) {
	var state State

	accounts := util.LoadAccounts()
	flag := false
	for _, ac := range accounts {
		if name == ac[0] {
			state.Name = name
			state.Address = ac[1]
			state.Key = ac[2]
			flag = true
			break
		}
	}
	if !flag {
		return nil, errors.New("このアカウント名は存在しません。")
	}
	return &state, nil
}
func GetAccount() *State {
	godotenv.Load()
	var state State
	state.Name = os.Getenv("NAME_BTC")
	state.Address = os.Getenv("ADDRESS_BTC")
	state.Key = os.Getenv("PRIVKEY_ENCRYPTED_BTC")
	return &state
}

func SaveConfig(st State) {
	curr, err := godotenv.Read(".env")
	if err != nil {
		curr = make(map[string]string)
	}
	curr["NAME_BTC"] = st.Name
	curr["ADDRESS_BTC"] = st.Address
	curr["PRIVKEY_ENCRYPTED_BTC"] = st.Key
	godotenv.Write(curr, ".env")
}

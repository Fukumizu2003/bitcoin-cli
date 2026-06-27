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
	state.Name = os.Getenv("NAME")
	state.Address = os.Getenv("ADDRESS")
	state.Key = os.Getenv("PRIVKEY_ENCRYPTED")
	return &state
}

func SaveConfig(st State) {
	curr, err := godotenv.Read(".env")
	if err != nil {
		curr = make(map[string]string)
	}
	curr["NAME"] = st.Name
	curr["ADDRESS"] = st.Address
	curr["PRIVKEY_ENCRYPTED"] = st.Key
	godotenv.Write(curr, ".env")
}

package config

import (
	"bitcoin-cli/internal/util"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type State struct {
	Name    string
	Address string
	Key     string
}

type Config struct{}

func ChangeMainAccount(address string) error {
	var state State
	util.MkdirOrNothing("ref")
	f, _ := os.ReadFile(filepath.Join("ref", "state.json"))
	json.Unmarshal(f, &state)

	accounts := util.LoadAccounts()
	flag := false
	for _, ac := range accounts {
		if address == ac[1] {
			state.Address = address
			state.Name = ac[0]
			state.Key = ac[2]
			flag = true
			break
		}
	}
	if !flag {
		return errors.New("このアドレスは登録されていません。")
	}

	stateSave, _ := json.MarshalIndent(state, "", "    ")
	os.WriteFile(filepath.Join("ref", "state.json"), stateSave, 0644)

	return nil
}

func GetMainAccount() State {
	var state State
	f, _ := os.ReadFile(filepath.Join("ref", "state.json"))
	json.Unmarshal(f, &state)
	return state
}

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

func Change_main_account(address string) error {
	var state State
	util.Mkdir_or_nothing("ref")
	f, _ := os.ReadFile(filepath.Join("ref", "state.json"))
	json.Unmarshal(f, &state)

	accounts := util.Load_accounts()
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

	state_save, _ := json.MarshalIndent(state, "", "    ")
	os.WriteFile(filepath.Join("ref", "state.json"), state_save, 0644)

	return nil
}

func Get_main_account() State {
	var state State
	f, _ := os.ReadFile(filepath.Join("ref", "state.json"))
	json.Unmarshal(f, &state)
	return state
}

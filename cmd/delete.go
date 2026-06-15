/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bitcoin-cli/internal/util"
	"bytes"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var delete_name string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if delete_name == "" {
			return fmt.Errorf("削除するアカウント名を-nフラグで指定してください。")
		}
		accounts := util.Load_accounts()
		destinations := util.Load_destinations()
		new_accounts := [][]string{}
		new_destinations := [][]string{}
		acflag := false
		deflag := false
		for _, ac := range accounts {
			if ac[0] != delete_name {
				new_accounts = append(new_accounts, ac)
			} else {
				acflag = true
			}
		}
		for _, ds := range destinations {
			if ds[0] != delete_name {
				new_destinations = append(new_destinations, ds)
			} else {
				deflag = true
			}
		}
		if !acflag && !deflag {
			return fmt.Errorf("アカウント名が存在しません。")
		}
		var buf bytes.Buffer
		writer := csv.NewWriter(&buf)
		if acflag {
			writer.WriteAll(new_accounts)
			os.WriteFile(util.Relative_to_absolute("ref", "keypair.csv"), buf.Bytes(), 0644)
		} else if deflag {
			writer.WriteAll(new_destinations)
			os.WriteFile(util.Relative_to_absolute("ref", "destinations.csv"), buf.Bytes(), 0644)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&delete_name, "name", "n", "", "削除するアカウントの名前を指定")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

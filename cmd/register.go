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

var address_only bool
var change_name bool
var set_address string
var set_name string
var from_name string
var to_name string
var set_password string

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if change_name {
			if from_name == "" && to_name == "" {
				fmt.Println("--from, --toフラグにより変更前の名前、変更後の名前を指定してください。")
				return
			}
			acs := util.Load_accounts()
			dss := util.Load_destinations()
			new_acs := [][]string{}
			new_dss := [][]string{}
			acflag := false
			deflag := false
			for _, ac := range acs {
				if ac[0] != from_name {
					new_acs = append(new_acs, ac)
				} else {
					new_acs = append(new_acs, []string{to_name, ac[1], ac[2]})
					acflag = true
				}
			}
			for _, ds := range dss {
				if ds[0] != from_name {
					new_dss = append(new_dss, ds)
				} else {
					new_dss = append(new_dss, []string{to_name, ds[1], ds[2]})
					deflag = true
				}
			}
			if !acflag && !deflag {
				fmt.Println("指定の名前のアカウントは存在しません。")
				return
			}
			var buf bytes.Buffer
			writer := csv.NewWriter(&buf)
			if acflag {
				writer.WriteAll(new_acs)
				os.WriteFile(util.Relative_to_absolute("ref", "keypair.csv"), buf.Bytes(), 0644)
			} else if deflag {
				writer.WriteAll(new_dss)
				os.WriteFile(util.Relative_to_absolute("ref", "destinations.csv"), buf.Bytes(), 0644)
			}
			return
		}
		if set_name == "" {
			fmt.Println("-nフラグによりアカウント名を指定してください。")
			return
		}
		acs := util.Load_accounts()
		dss := util.Load_destinations()
		if !util.Check_name(acs, dss, set_name) {
			fmt.Println("この名前は既に登録されています。同じ名前のアカウントを重複して作ることはできません。")
			return
		}
		if !address_only {
			if set_password == "" {
				fmt.Println("-pフラグによりパスワードを設定してください。空白を含む場合は、二重引用符\"\"で囲んでください。")
			}
			privkey, pubkey := util.New_keypair()
			address := util.Pubkey_to_address_b32(pubkey.SerializeCompressed())
			privkey_cr := util.Aes_encrypt(privkey.Serialize(), []byte(set_password))
			util.Save_keypair(set_name, address, privkey_cr)
		} else {
			if set_address != "" {
				util.Save_address(set_name, set_address)
			} else {
				fmt.Println("-aフラグによりアドレスを指定してください。")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().BoolVarP(&address_only, "addressonly", "o", false, "")
	registerCmd.Flags().BoolVarP(&change_name, "change", "c", false, "")
	registerCmd.Flags().StringVarP(&set_address, "address", "a", "", "")
	registerCmd.Flags().StringVarP(&set_name, "name", "n", "", "")
	registerCmd.Flags().StringVarP(&from_name, "from", "f", "", "")
	registerCmd.Flags().StringVarP(&to_name, "to", "t", "", "")
	registerCmd.Flags().StringVarP(&set_password, "password", "p", "", "")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

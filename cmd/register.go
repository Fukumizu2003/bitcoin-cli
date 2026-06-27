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

var addressOnly bool
var changeName bool
var setAddress string
var setName string
var fromName string
var toName string
var setPassword string

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if changeName {
			if fromName == "" && toName == "" {
				fmt.Println("--from, --toフラグにより変更前の名前、変更後の名前を指定してください。")
				return
			}
			acs := util.LoadAccounts()
			dss := util.LoadDestinations()
			newAcs := [][]string{}
			newDss := [][]string{}
			acflag := false
			deflag := false
			for _, ac := range acs {
				if ac[0] != fromName {
					newAcs = append(newAcs, ac)
				} else {
					newAcs = append(newAcs, []string{toName, ac[1], ac[2]})
					acflag = true
				}
			}
			for _, ds := range dss {
				if ds[0] != fromName {
					newDss = append(newDss, ds)
				} else {
					newDss = append(newDss, []string{toName, ds[1], ds[2]})
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
				writer.WriteAll(newAcs)
				os.WriteFile(util.RelativeToAbsolute("ref", "keypair.csv"), buf.Bytes(), 0644)
			} else if deflag {
				writer.WriteAll(newDss)
				os.WriteFile(util.RelativeToAbsolute("ref", "destinations.csv"), buf.Bytes(), 0644)
			}
			return
		}
		if setName == "" {
			fmt.Println("-nフラグによりアカウント名を指定してください。")
			return
		}
		acs := util.LoadAccounts()
		dss := util.LoadDestinations()
		if !util.CheckName(acs, dss, setName) {
			fmt.Println("この名前は既に登録されています。同じ名前のアカウントを重複して作ることはできません。")
			return
		}
		if !addressOnly {
			if setPassword == "" {
				fmt.Println("-pフラグによりパスワードを設定してください。空白を含む場合は、二重引用符\"\"で囲んでください。")
			}
			privkey, pubkey := util.NewKeypair()
			address := util.PubkeyToAddressB32(pubkey.SerializeCompressed())
			privkeyCr := util.AesEncrypt(privkey.Serialize(), []byte(setPassword))
			util.SaveKeypair(setName, address, privkeyCr)
		} else {
			if setAddress != "" {
				if !util.IsValidAddress(setAddress) {
					fmt.Println("指定のアドレスは無効、もしくは未対応の形式です。")
					return
				}
				util.SaveAddress(setName, setAddress)
			} else {
				fmt.Println("-aフラグによりアドレスを指定してください。")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().BoolVarP(&addressOnly, "addressonly", "o", false, "")
	registerCmd.Flags().BoolVarP(&changeName, "change", "c", false, "")
	registerCmd.Flags().StringVarP(&setAddress, "address", "a", "", "")
	registerCmd.Flags().StringVarP(&setName, "name", "n", "", "")
	registerCmd.Flags().StringVarP(&fromName, "from", "f", "", "")
	registerCmd.Flags().StringVarP(&toName, "to", "t", "", "")
	registerCmd.Flags().StringVarP(&setPassword, "password", "p", "", "")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

package config

import (
	"bufio"
	"fmt"
	"github.com/bsthun/gut"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
	"log"
	"os"
	"strawhouse-command/common"
	"strings"
)

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a config key and value",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if err := gut.Validator.Var(name, "oneof=server secure key"); err != nil {
			log.Fatalf("name is required one of 'server', 'secure', or 'key'")
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("enter %s value: ", name)
		value, _ := reader.ReadString('\n')
		value = strings.Replace(value, "\n", "", -1)
		if err := gut.Validator.Var(value, "required"); err != nil {
			log.Fatalf("value is required for key")
		}

		viper.Set(name, value)
		if name == "key" {
			err := keyring.Set(common.KeyringService, common.KeyringUser, value)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("key set in keyring")
		} else {
			if err := viper.WriteConfig(); err != nil {
				log.Fatalf("unable to write to config file: %v", err)
			}
			fmt.Printf("%s configuration set to '%s'\n", name, value)
		}
	},
}

func init() {
	setCmd.Flags().String("name", "", "Config to set. One of 'server', 'secure', or 'key'")
	_ = setCmd.MarkFlagRequired("name")
	Cmd.AddCommand(setCmd)
}

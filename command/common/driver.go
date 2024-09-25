package common

import (
	strawhouse "driver"
	"github.com/spf13/viper"
	"log"
)

var Driver *strawhouse.Driver

func InitDriver() {
	key := ConfigVaultKeyGet()
	server := viper.Get("server")
	if server == nil {
		log.Fatalf("server is required, please use 'strawc config set --name server --value <server>'")
	}
	Driver = strawhouse.New(key, server.(string))
}

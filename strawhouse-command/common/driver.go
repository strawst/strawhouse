package common

import (
	"github.com/spf13/viper"
	"github.com/strawstacks/strawhouse/strawhouse-driver"
	"log"
)

var Driver *strawhouse.Driver

func InitDriver() {
	key := ConfigVaultKeyGet()
	server := viper.Get("server")
	secure := viper.Get("secure")
	if server == nil {
		log.Fatalf("server is required, please use 'strawc config set --name server'")
	}
	if secure == nil {
		log.Fatalf("secure is required, please use 'strawc config set --name secure' with value of <y/n>'")
	}
	Driver = strawhouse.New(key, server.(string), &strawhouse.Option{
		Secure: secure.(string) == "y",
	})
}

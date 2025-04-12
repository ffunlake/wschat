package config

import (
	"os"
	"strings"
	"log"
	"fmt"
	"github.com/spf13/viper"
)


func LoadConfig(configFile string) error {
	viper.SetConfigFile(configFile)
	content, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}
	//Replace environment variables
	return viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
}



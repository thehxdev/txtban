package config

import (
	"github.com/spf13/viper"
)

type TbConfig struct {
	Server struct {
		Port    int
		Address string
	}

	Database struct {
		Path string
	}

	Limits struct {
		MinPasswordLen int
		MaxTxtIdLen    int
		MaxTxtNameLen  int
	}
}

func SetupViper(c *TbConfig, configPath string) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	viper.AddConfigPath("/etc/txtban/")
	viper.AddConfigPath(".")

	if configPath != "" {
		viper.SetConfigFile(configPath)
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}

	err = viper.Unmarshal(c)
	if err != nil {
		panic(err.Error())
	}
}

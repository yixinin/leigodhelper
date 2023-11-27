package helper

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Games    string `mapstructure:"games"`
}

var config Config

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if configPath != "" {
		if strings.Contains(configPath, ".") {
			viper.SetConfigFile(configPath)
		} else {
			viper.SetConfigName(configPath)
		}
	}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	config.Games = strings.ReplaceAll(config.Games, "，", ",")
	return &config, nil
}

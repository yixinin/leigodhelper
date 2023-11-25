package helper

import (
	"fmt"
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
		Logger.Printf("load config file:%s\n", configPath)
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

	if config.Username == "" || config.Password == "" {
		return nil, fmt.Errorf("雷神加速器助手已退出，请正确填写雷神加速器账号密码后重启")
	}

	// if _, ok := config.StartWith["steam"]; !ok {
	// 	config.StartWith["steam"] = "C:\\Program Files (x86)\\Steam\\steam.exe"
	// }
	// if _, ok := config.StartWith["leigod"]; !ok {
	// 	config.StartWith["leigod"] = "C:\\Program Files (x86)\\LeiGod_Acc\\leigod.exe"
	// }

	config.Games = strings.ReplaceAll(config.Games, "，", ",")
	Logger.Println("start watch games:", config.Games)
	return &config, nil
}

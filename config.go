package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Username  string            `mapstructure:"username"`
	Password  string            `mapstructure:"password"`
	Games     string            `mapstructure:"games"`
	StartWith map[string]string `mapstructure:"start_with"` // 进程名:启动路径
}

var config Config

func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}

	// if _, ok := config.StartWith["steam"]; !ok {
	// 	config.StartWith["steam"] = "C:\\Program Files (x86)\\Steam\\steam.exe"
	// }
	// if _, ok := config.StartWith["leigod"]; !ok {
	// 	config.StartWith["leigod"] = "C:\\Program Files (x86)\\LeiGod_Acc\\leigod.exe"
	// }

	strings.ReplaceAll(config.Games, "，", ",")
	var games = strings.Split(config.Games, ",")
	watchs = make(map[string]struct{}, len(games))
	for _, game := range games {
		game = strings.TrimSpace(game)
		watchs[game] = struct{}{}
	}
	return nil
}

var Logger *log.Logger

func init() {
	cmd := exec.Command("powershell", "rm", "leigodhelper.log.*")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	os.Rename("leigodhelper.log", fmt.Sprintf("leigodhelper.log.%s", time.Now().Format("0102150405")))
	f, err := os.Create("leigodhelper.log")
	if err != nil {
		log.Println(err)
		return
	}
	Logger = log.New(f, "", log.Llongfile)
}

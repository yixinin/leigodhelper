package main

import (
	"flag"
	"fmt"
	"leigodhelper/helper"
)

var configPath string

func main() {
	defer func() {
		if r := recover(); r != nil {
			var msg = fmt.Sprintf("雷神加速器助手异常退出\nerror: %+v", r)
			helper.Notify(msg)
			helper.Logger.Println(msg)
		}
		helper.Logger.Println("助手已正常关闭")
	}()
	flag.Parse()
	var args = flag.CommandLine.Args()
	if len(args) > 0 {
		configPath = args[0]
	}
	if err := helper.LoadConfig(configPath); err != nil {
		helper.Notify(err.Error())
		return
	}
	helper.LoadToken()
	helper.Run()
}

package main

import (
	"fmt"
	"leigodhelper/helper"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			var msg = fmt.Sprintf("雷神加速器助手异常退出\nerror: %+v", r)
			helper.Notify(msg)
			helper.Logger.Println(msg)
		}
		helper.Logger.Println("助手已正常关闭")
	}()
	if err := helper.LoadConfig(); err != nil {
		helper.Notify(err.Error())
		return
	}
	helper.LoadToken()
	helper.Run()
}

package main

import (
	"github.com/go-toast/toast"
)

func notify(msg string) {
	notification := toast.Notification{
		AppID:   "Microsoft.Windows.Shell.RunDialog",
		Title:   "加速器助手",
		Message: msg,
	}
	err := notification.Push()
	if err != nil {
		Logger.Fatalln(err)
	}
}

package main

import (
	"os/exec"
	"syscall"
)

func runApps(apps ...string) error {
	// 启动
	for _, app := range apps {
		cmd := exec.Command("cmd", "/c", app)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err := cmd.Start()
		if err != nil {
			Logger.Println(err)
			return err
		}
		Logger.Printf("start %s\n", app)
	}
	return nil
}

package helper

import (
	"os/exec"
	"syscall"
	"time"
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
		Logger.Println("start", app, time.Now())
	}
	return nil
}

package helper

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var Logger *Log

func init() {
	cmd := exec.Command("powershell", "rm", "logs/leigodhelper.log.*")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	os.Rename("logs/leigodhelper.log", fmt.Sprintf("logs/leigodhelper.log.%s", time.Now().Format("0102150405")))
	os.Mkdir("logs", os.ModeDir)
	f, err := os.Create("logs/leigodhelper.log")
	if err != nil {
		log.Println(err)
		return
	}
	logger := log.New(f, "", log.Llongfile)

	Logger = &Log{
		logger: logger,
	}
}

type Log struct {
	logger *log.Logger
}

func (l *Log) Println(args ...interface{}) {
	args = append([]interface{}{time.Now()}, args...)
	l.logger.Println(args...)
}
func (l *Log) Printf(f string, args ...interface{}) {
	args = append([]interface{}{time.Now()}, args...)
	l.logger.Printf("%s "+f, args...)
}

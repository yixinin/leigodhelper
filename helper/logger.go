package helper

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

var Logger *Log

func init() {
	var dir = `C:\Program Files\LeigodHelper`
	cmd := exec.Command("powershell", "rm", filepath.Join(dir, "logs/leigodhelper.log.*"))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()

	os.Rename(filepath.Join(dir, "logs/leigodhelper.log"), fmt.Sprintf(filepath.Join(dir, "logs/leigodhelper.log.%s"), time.Now().Format("0102150405")))
	os.Mkdir(filepath.Join(dir, "logs"), os.ModeDir)
	f, err := os.OpenFile(filepath.Join(dir, "logs/leigodhelper.log"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	logger := log.New(f, "", log.Llongfile)

	Logger = &Log{
		f:      f,
		logger: logger,
	}
}

type Log struct {
	f      *os.File
	logger *log.Logger
}

func (l *Log) Close() {
	if l != nil {
		l.f.Close()
	}
}

func (l *Log) Println(args ...interface{}) {
	args = append([]interface{}{time.Now()}, args...)
	l.logger.Println(args...)
}
func (l *Log) Printf(f string, args ...interface{}) {
	args = append([]interface{}{time.Now()}, args...)
	l.logger.Printf("%s "+f, args...)
}

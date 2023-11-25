package main

import (
	"context"
	"fmt"
	"io"
	"leigodhelper/helper"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/kardianos/service"
)

const InstallPath = `C:\Program Files\LeigodHelper`

func main() {
	defer func() {
		if r := recover(); r != nil {
			var msg = fmt.Sprintf("leigod helper exited with error:\nerror: %+v", r)
			helper.Notify(msg)
			helper.Logger.Println(msg)
			return
		}
	}()
	svcConfig := &service.Config{
		Name:             "LeigodHelper",
		DisplayName:      "LeigodHelper",
		Description:      "Leigod game helper service",
		WorkingDirectory: `C:\Program Files\LeigodHelper`,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		helper.Notify(err.Error())
		return
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			os.MkdirAll(InstallPath, 0644)
			src, err := os.Open(os.Args[0])
			if err != nil {
				panic(err)
			}
			dst, err := os.Create(InstallPath + "/leigodhelper.exe")
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(dst, src)
			if err != nil {
				panic(err)
			}
			// start service
			s.Install()
			log.Println("Install service success!")
			err = exec.Command("net", "start", "LeigodHelper").Run()
			if err != nil {
				log.Println("start service error: ", err)
			}
			return
		}

		if os.Args[1] == "remove" {
			s.Uninstall()
			err := os.RemoveAll(InstallPath)
			if err != nil {
				log.Println("remove files error: ", err)
			}
			err = exec.Command("net", "stop", "LeigodHelper").Run()
			if err != nil {
				log.Println("stop service error: ", err)
			}
			log.Println("Remove service success!")
			return
		}
	}

	// Call svc.Run to start your program/service.
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

// program implements svc.Service
type program struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func (p *program) Start(srv service.Service) error {
	c, err := helper.LoadConfig(`C:\Program Files\LeigodHelper\config.toml`)
	if err != nil {
		helper.Notify(err.Error())
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	var h = helper.NewHelper(c)
	p.wg.Add(1)
	go func() {
		h.Run(ctx)
		p.wg.Done()
	}()
	helper.Logger.Println("Starting...")
	return nil
}

func (p *program) Stop(srv service.Service) error {
	// The Stop method is invoked by stopping the Windows service, or by pressing Ctrl+C on the console.
	// This method may block, but it's a good idea to finish quickly or your process may be killed by
	// Windows during a shutdown/reboot. As a general rule you shouldn't rely on graceful shutdown.
	if p.cancel != nil {
		p.cancel()
	}
	helper.Logger.Println("Waiting graceful stop...")
	p.wg.Wait()
	helper.Logger.Println("Stopped")
	return nil
}

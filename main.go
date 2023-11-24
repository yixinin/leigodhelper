package main

import (
	"context"
	"fmt"
	"leigodhelper/helper"
	"log"
	"os"
	"sync"

	"github.com/kardianos/service"
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
			s.Install()
			log.Println("服务安装成功")
			return
		}

		if os.Args[1] == "remove" {
			s.Uninstall()
			log.Println("服务卸载成功")
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
	quit   chan struct{}
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
	go h.Run(ctx)
	// The Start method must not block, or Windows may assume your service failed
	// to start. Launch a Goroutine here to do something interesting/blocking.

	p.quit = make(chan struct{})

	p.wg.Add(1)
	go func() {
		log.Println("Starting...")
		<-p.quit
		log.Println("Quit signal received...")
		p.wg.Done()
	}()

	return nil
}

func (p *program) Stop(srv service.Service) error {
	// The Stop method is invoked by stopping the Windows service, or by pressing Ctrl+C on the console.
	// This method may block, but it's a good idea to finish quickly or your process may be killed by
	// Windows during a shutdown/reboot. As a general rule you shouldn't rely on graceful shutdown.
	if p.cancel != nil {
		p.cancel()
	}
	log.Println("Stopping...")
	close(p.quit)
	p.wg.Wait()
	log.Println("Stopped.")
	return nil
}

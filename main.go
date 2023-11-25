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
			s.Install()
			log.Println("Install service success!")
			return
		}

		if os.Args[1] == "remove" {
			s.Uninstall()
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

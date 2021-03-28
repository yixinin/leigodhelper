package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	watchs map[string]struct{}
	leigod = "leigod"

	lastTryPauseTime int64
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			notify(fmt.Sprintf("系统错误，关闭助手, main\nerr: %+v", r))
		}
		Logger.Println("defer main")
	}()
	if err := loadConfig(); err != nil {
		notify(err.Error())
	}
	LoadToken()
	run()
}

func run() {
	var exitCh = make(chan string)
	sysCh := make(chan os.Signal)
	signal.Notify(sysCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	// 启动加速器和steam
	ctx, cancelWait := context.WithTimeout(context.Background(), 60*time.Second)

	newStarted := waitStart(ctx)
	if newStarted { // 刚打开雷神加速器 不检测暂停
		lastTryPauseTime = time.Now().Unix()
	}

	cancelWait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	Logger.Println("start check", time.Now())
	go check(ctx, exitCh)

	for {
		select {
		case <-sysCh:
			tryPauseLeigod()
			Logger.Printf("加速器助手已被强制退出")
			return
		case msg := <-exitCh:
			Logger.Printf("加速器助手已退出, msg: %s", msg)
			tryPauseLeigod()
			return
		}
	}
}

func waitStart(ctx context.Context) bool {
	var tryRun = false
	var startMap = make(map[string]bool, len(config.StartWith))
	for k := range config.StartWith {
		startMap[k] = false
	}
	for {
		select {
		case <-ctx.Done(): //超时退出
			return tryRun
		default:
		}
		var apps = make([]string, 0, len(config.StartWith))
		// 尝试启动雷神加速器和steam
		var pMap = getProcess()

		// 检测加速器是否运行中
		for pName, started := range startMap {
			if started {
				continue
			}
			if pNum, ok := pMap[pName]; ok && pNum > 0 {
				startMap[pName] = true
			} else {
				apps = append(apps, config.StartWith[pName])
			}
		}
		// 已启动
		if len(apps) == 0 {
			return tryRun
		}

		if !tryRun {
			tryRun = true
			if err := runApps(apps...); err != nil {
				Logger.Println(err)
				return false
			}
		}
		time.Sleep(time.Second)
	}
}

func check(ctx context.Context, exitCh chan string) {
	defer func() {
		if r := recover(); r != nil {
			exitCh <- fmt.Sprintf("系统错误，关闭助手\n err: %+v", recover())
		}
		Logger.Println("defer check")
	}()
	for {

		select {
		case <-ctx.Done():
			return
		default:
		}
		leigodOK, gameOK := hasGameRunning()
		if !leigodOK {
			exitCh <- "" // 加速器退出时 助手也退出
		}

		if gameOK { // 启动游戏时 重置检测状态
			lastTryPauseTime = -1
			time.Sleep(5 * time.Second)
			continue
		}

		if !gameOK && time.Now().Unix()-lastTryPauseTime > 60*60 { // 如果近期检测过 不再频繁检测
			Logger.Printf("lastTryPauseTime:%d\n", lastTryPauseTime)
			if err := tryPauseLeigod(); err != nil {
				notify("加速器助手已退出：" + err.Error())
				exitCh <- ""
				return
			}
		}
		time.Sleep(2 * time.Second)
	}
}

// 检测游戏是否运行中
func hasGameRunning() (leigodOK, gameOK bool) {
	var ps = getProcess()
	if _, leigodOK = ps["leigod"]; !leigodOK {
		return
	}
	for k := range watchs {
		if _, ok := ps[k]; ok {
			gameOK = true
			break
		}
	}
	return
}

func tryPauseLeigod() error {
	lastTryPauseTime = time.Now().Unix()
	if !Token.IsValid() {
		token, expire, err := Login(config.Username, config.Password)
		if err != nil {
			return err
		}
		Token = &LeigodToken{
			AccountToken: token,
			ExpireTime:   expire,
		}
		SaveToken()
	}
	// 获取暂停状态
	paused, err := IsPause(Token.AccountToken)
	if err != nil {
		return err
	}
	if !paused {
		err := Pause(Token.AccountToken)
		if err != nil {
			return err
		}
		notify("雷神加速器助手检测到当前没有游戏运行，已暂停时长。")
		return nil
	}
	Logger.Println("已是暂停状态，无需暂停")
	return nil
}

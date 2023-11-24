package helper

import (
	"context"
	"fmt"
	"strings"
	"time"
)

var (
	lastTryPauseTime int64
)

type Helper struct {
	api     LeigodApi
	games   map[string]bool
	timeout time.Duration

	GameStatus   Status
	LeigodStatus Status

	tm    *time.Timer
	tmSet bool
}

type Status string

const (
	Unknown Status = ""
	Stop    Status = "stop"
	Running Status = "running"
)

func NewHelper(c *Config) Helper {
	var games = make(map[string]bool)
	for _, v := range strings.Split(c.Games, ",") {
		games[strings.TrimSpace(v)] = true
	}

	h := Helper{
		api: LeigodApi{
			username: c.Username,
			password: c.Password,
		},
		games:   games,
		timeout: time.Duration(c.Timeout) * time.Second,
		tm:      time.NewTimer(time.Second),
	}
	h.tm.Stop()
	return h
}

func (h *Helper) Run(ctx context.Context) {
	var exitCh = make(chan string)
	Logger.Println("start check")
	go h.loop(ctx, exitCh)

	for {
		select {
		case <-ctx.Done():
			h.Pause()
			Logger.Printf("加速器助手已被强制退出")
			return
		case msg := <-exitCh:
			Logger.Printf("加速器助手已退出, msg: %s", msg)
			h.Pause()
			return
		case <-h.tm.C:
			h.Pause()
			h.tmSet = false
		}
	}
}

func (h *Helper) Update(leigodOK, gameOK bool) {
	h.GameStatus = Stop
	if gameOK {
		h.GameStatus = Running
		if h.tmSet {
			Logger.Println("game running, pause canced!")
			h.tm.Stop()
			h.tmSet = false
			return
		}
	}
	switch h.LeigodStatus {
	case Running, Unknown:
		if !leigodOK || !gameOK {
			if !h.tmSet {
				Logger.Println("will pause 3 minute later...")
				if h.LeigodStatus == Running {
					Notify("雷神加速器助手检测到当前没有游戏运行，3分钟后停时长。")
				}
				h.tm.Reset(3 * time.Minute)
				h.tmSet = true
			}
		}
	}
}

func (h *Helper) loop(ctx context.Context, exitCh chan string) {
	defer func() {
		if r := recover(); r != nil {
			exitCh <- fmt.Sprintf("系统错误，关闭助手\n err: %+v", recover())
		}
		Logger.Println("defer check")
	}()
	tk := time.NewTicker(10 * time.Second)
	defer tk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			leigodOk, gameOK := hasGameRunning("leigod", h.games)
			h.Update(leigodOk, gameOK)
		}
	}
}

// 检测游戏是否运行中
func hasGameRunning(name string, games map[string]bool) (leigodOK, gameOK bool) {
	var ps = getProcess()
	if _, leigodOK = ps[name]; !leigodOK {
		return
	}
	for k := range games {
		if _, ok := ps[k]; ok {
			gameOK = true
			break
		}
	}
	return
}

func (h *Helper) Pause(finnal ...bool) error {
	lastTryPauseTime = time.Now().Unix()
	// 获取暂停状态
	paused, err := h.api.IsPause()
	if err == ErrorNotLogin {
		if err := h.Relogin(); err != nil {
			return err
		}
		if len(finnal) > 0 {
			return nil
		}
		return h.Pause(true)
	} else {
		if err != nil {
			return err
		}
	}
	if !paused {
		h.LeigodStatus = Running
		err := h.api.Pause()
		if err == ErrorNotLogin {
			if err := h.Relogin(); err != nil {
				return err
			}
			if len(finnal) > 0 {
				return nil
			}
			return h.Pause(true)
		} else {
			if err != nil {
				return err
			}
		}

		Notify("雷神加速器助手检测到当前没有游戏运行，已暂停时长。")
		h.LeigodStatus = Stop
		return nil
	}
	h.LeigodStatus = Stop
	Logger.Println("已是暂停状态，无需暂停")
	return nil
}

func (h *Helper) Relogin() error {
	err := h.api.Login()
	if err != nil {
		return err
	}
	return nil
}

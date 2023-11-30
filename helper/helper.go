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
	api   LeigodApi
	games map[string]bool

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
		name := NormalizeProcessName(strings.TrimSpace(v))
		games[name] = true
	}

	h := Helper{
		api: LeigodApi{
			username: c.Username,
			password: c.Password,
		},
		games: games,
		tm:    time.NewTimer(time.Second),
	}
	h.tm.Stop()
	return h
}

func (h *Helper) Run(ctx context.Context) {
	var exitCh = make(chan string)
	Logger.Println("start watch games:", config.Games)
	go h.loop(ctx, exitCh)
	for {
		select {
		case <-ctx.Done():
			Logger.Printf("recv exit signal, try pause")
			h.Pause()
			return
		case msg := <-exitCh:
			Logger.Printf("check paniced, try pause", msg)
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
	if leigodOK {
		if !gameOK && !h.tmSet {
			Logger.Println("will pause 10minute later...")
			if h.LeigodStatus == Running {
				Notify("no game running, will pause 10minutes later.")
			}
			h.tm.Reset(10 * time.Minute)
			h.tmSet = true
		}
		return
	} else {
		if h.LeigodStatus == Running {
			if !gameOK && !h.tmSet {
				Logger.Println("will pause 10minute later...")
				if h.LeigodStatus == Running {
					Notify("no game running, will pause 10minutes later.")
				}
				h.tm.Reset(10 * time.Minute)
				h.tmSet = true
			}
		}
	}
}

func (h *Helper) loop(ctx context.Context, exitCh chan string) {
	defer func() {
		if r := recover(); r != nil {
			exitCh <- fmt.Sprintf("paniced ! \n err: %+v", recover())
		}
	}()
	tk := time.NewTicker(10 * time.Second)
	defer tk.Stop()

	ttk := time.NewTicker(time.Hour)
	defer ttk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			c, _ := LoadConfig(`C:\Program Files\LeigodHelper\config.toml`)
			if c != nil {
				var games = make(map[string]bool)
				var adds = make([]string, 0)
				var dels = make([]string, 0)
				for _, v := range strings.Split(c.Games, ",") {
					name := NormalizeProcessName(strings.TrimSpace(v))
					games[name] = true

					if _, ok := h.games[name]; !ok {
						adds = append(adds, name)
					}
				}
				for k := range h.games {
					if _, ok := games[k]; !ok {
						dels = append(dels, k)
					}
				}
				if len(adds) != 0 {
					Logger.Println("add wathcing games:", dels)
				}
				if len(dels) != 0 {
					Logger.Println("del wathcing games:", dels)
				}
				if c.Password != "" && c.Password != h.api.password {
					Logger.Println("update password")
					h.api.password = c.Password
				}
				if c.Username != "" && c.Username != h.api.username {
					Logger.Println("update username to", c.Username)
					h.api.username = c.Username
				}

				h.games = games
			}
			leigodOk, gameOK := hasGameRunning("leigod", h.games)
			h.Update(leigodOk, gameOK)
		case <-ttk.C:
			// force check
			_, gameOK := hasGameRunning("leigod", h.games)
			if !gameOK {
				h.Pause()
			}
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

		Notify("no game running, stop timing.")
		h.LeigodStatus = Stop
		return nil
	}
	h.LeigodStatus = Stop
	Logger.Println("already paused!, nothing todo.")
	return nil
}

func (h *Helper) Relogin() error {
	err := h.api.Login()
	if err != nil {
		return err
	}
	return nil
}

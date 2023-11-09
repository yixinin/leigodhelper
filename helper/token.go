package helper

import (
	"encoding/json"
	"io/fs"
	"os"
	"time"
)

type LeigodToken struct {
	AccountToken string
	ExpireTime   time.Time
}

var Token *LeigodToken

var filename = "cache/token.cache"

func (a LeigodApi) SaveToken() {
	buf, err := json.Marshal(a.Token)
	if err != nil {
		Logger.Println(err)
		return
	}
	os.Mkdir("cache", os.ModeDir)
	err = os.WriteFile(filename, buf, fs.ModePerm)
	if err != nil {
		Logger.Println(err)
	}
}

func (a LeigodApi) LoadToken() {
	buf, _ := os.ReadFile(filename)
	json.Unmarshal(buf, &a.Token)
}

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

func (t *LeigodToken) IsValid() bool {
	if t == nil {
		return false
	}
	return t.ExpireTime.After(time.Now()) && t.AccountToken != ""
}

var Token *LeigodToken

var filename = "cache/token.cache"

func SaveToken() {
	buf, err := json.Marshal(Token)
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

func LoadToken() {
	if Token == nil {
		Token = new(LeigodToken)
	}
	buf, _ := os.ReadFile(filename)
	json.Unmarshal(buf, Token)
}

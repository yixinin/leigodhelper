package main

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

var filename = "token.cache"

func SaveToken() {
	buf, err := json.Marshal(Token)
	if err != nil {
		Logger.Println(err)
		return
	}
	err = os.WriteFile(filename, buf, fs.ModePerm)
	if err != nil {
		Logger.Println(err)
	}
}

func LoadToken() {
	if Token == nil {
		Token = new(LeigodToken)
	}
	buf, err := os.ReadFile(filename)
	if err != nil {
		Logger.Println(err)
		return
	}
	err = json.Unmarshal(buf, Token)
	if err != nil {
		Logger.Println(err)
	}
}

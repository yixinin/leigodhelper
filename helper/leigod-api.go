package helper

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	CodeNotLogin = 400006
	TimeLayout   = "2006-01-02 15:04:05"

	UserAgent = "Mozilla/5.0 (Linux; Android 9; MIX 2 Build/PKQ1.190118.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045521 Mobile Safari/537.36 MMWEBID/2679 MicroMessenger/8.0.1.1841(0x2800015D) Process/tools WeChat/arm64 Weixin NetType/WIFI Language/zh_CN ABI/arm64"
)

var ErrorNotLogin = errors.New("not login")

type Ack struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type LoginReq struct {
	Lang        string `json:"lang"`
	CountryCode int    `json:"country_code"`
	Code        string `json:"code"`
	SrcChannel  string `json:"src_channel"`
	UserType    string `json:"user_type"`
	Password    string `json:"password"`
	Username    string `json:"username"`
}

type LoginAck struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		LoginInfo struct {
			AccountToken string `json:"account_token"`
			ExpiryTime   string `json:"expiry_time"`
		} `json:"login_info"`
	} `json:"data"`
}

type LeigodApi struct {
	username, password string
	Token              LeigodToken
}

func NewApi(username, password string) LeigodApi {
	api := LeigodApi{
		username: username,
		password: password,
	}
	api.LoadToken()
	return api
}

func (a LeigodApi) GetToken() string {
	if a.Token.ExpireTime.Before(time.Now()) || a.Token.AccountToken == "" {
		if err := a.Login(); err != nil {
			Logger.Println(err)
			return ""
		}
	}
	return a.Token.AccountToken
}

func (a LeigodApi) Login() error {
	var url = "https://webapi.nn.com/wap/login/bind"
	var req = LoginReq{
		Lang:        "zh_CN",
		CountryCode: 86,
		SrcChannel:  "guanwang",
		UserType:    "0",
		Username:    a.username,
		Password:    MD5(a.password),
	}
	var ack LoginAck
	err := a.leigodHttpPost(url, req, &ack)
	if err != nil {
		return err
	}
	if ack.Code != 0 {
		return fmt.Errorf("%s", ack.Message)
	}
	expireTime, _ := time.ParseInLocation(TimeLayout, ack.Data.LoginInfo.ExpiryTime, time.Local)
	a.Token.AccountToken = ack.Data.LoginInfo.AccountToken
	a.Token.ExpireTime = expireTime
	a.SaveToken()
	return nil
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

type PauseReq struct {
	AccountToken string `json:"account_token"`
	Lang         string `json:"lang"`
}

type PauseAck struct {
	Code    int           `json:"code"`
	Data    []interface{} `json:"data"`
	Message string        `json:"msg"`
}

func (a LeigodApi) Pause() error {
	var url = "https://webapi.nn.com/api/user/pause"
	var req = PauseReq{
		Lang:         "zh_CN",
		AccountToken: a.GetToken(),
	}
	var ack PauseAck
	err := a.leigodHttpPost(url, req, &ack)
	if err != nil {
		return err
	}
	if ack.Code != 0 {
		return fmt.Errorf("code:%d, msg:%s, data:%+v", ack.Code, ack.Message, ack.Data)
	}
	return nil
}

type UserInfoReq struct {
	AccountToken string `json:"account_token"`
	Lang         string `json:"lang"`
}

type UserInfoAck struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		PauseStatusId int `json:"pause_status_id"` // 0表示未暂停
	} `json:"data"`
}

func (a LeigodApi) IsPause() (bool, error) {
	var url = "https://webapi.nn.com/api/user/info"
	var req = UserInfoReq{
		Lang:         "zh_CN",
		AccountToken: a.GetToken(),
	}
	var ack UserInfoAck
	err := a.leigodHttpPost(url, req, &ack)
	if err != nil {
		return false, err
	}
	if ack.Code != 0 {
		return false, fmt.Errorf("code:%d, msg:%s", ack.Code, ack.Message)
	}
	return ack.Data.PauseStatusId != 0, nil
}

func (a LeigodApi) leigodHttpPost(url string, req, ackBody interface{}) error {
	reqBuf, err := json.Marshal(req)
	if err != nil {
		return err
	}
	r, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBuf))
	if err != nil {
		return err
	}
	a.setHeader(r)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// Logger.Printf("http post: url=%s, req=%+v, ack=%s \n", url, req, string(buf))
	var ack Ack
	err = json.Unmarshal(buf, &ack)
	if err != nil {
		return err
	}
	if ack.Code == 0 {
		err = json.Unmarshal(buf, ackBody)
		Logger.Printf("http post: url=%s, req=%s, ack=%+v \n", url, string(reqBuf), ackBody)
		return err
	}
	if ack.Code == CodeNotLogin {
		return ErrorNotLogin
	}
	Logger.Printf("http post: url=%s, req=%s, ack=%+v \n", url, string(reqBuf), ack)
	return fmt.Errorf("code:%d, msg:%s", ack.Code, ack.Message)
}

func (a LeigodApi) setHeader(r *http.Request) {
	var headers = map[string]string{
		"Origin":           "https://jiasu.nn.com/",
		"User-Agent":       UserAgent,
		"Sec-Fetch-Mode":   "cors",
		"X-Requested-With": "com.tencent.mm",
		"Sec-Fetch-Site":   "same-site",
		"Referer":          "https://jiasu.nn.com/m/mpause.html?region_code=1&language=zh_CN&platform=4",
	}
	for k, v := range headers {
		r.Header.Set(k, v)
	}
}

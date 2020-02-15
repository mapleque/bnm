package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WxSDK struct {
	AppId  string
	Secret string
	log    Logger
}

type WxAuthenticate struct {
	OpenId  string
	UnionId string
}

func (wx *WxSDK) Jscode2session(code string) (*WxAuthenticate, error) {
	// @see document online:
	// https://developers.weixin.qq.com/miniprogram/dev/api/open-api/login/code2Session.html
	urlTpl := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	url := fmt.Sprintf(urlTpl, wx.AppId, wx.Secret, code)

	ret, err := wx.HttpGet(url)
	if err != nil {
		return nil, err
	}
	auth := &WxAuthenticate{}
	if err := json.Unmarshal(ret, auth); err != nil || auth.OpenId == "" {
		return nil, fmt.Errorf("%s", string(ret))
	}
	return auth, nil
}

func (wx *WxSDK) HttpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	wx.log.Trace("HTTP GET", url, string(body))
	return body, nil
}

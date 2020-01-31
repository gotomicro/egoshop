package conf

import (
	"errors"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

var (
	Conf config // holds the global app config.
)

type config struct {
	// 应用配置
	App app
	Oss oss
}

type app struct {
	Name       string `toml:"name"`
	Wechat     wechat `toml:"wechat"`
	WechatPay  wechatPay
	WechatOpen wechatOpen
	CdnName    string
	File       string `toml:"file"`
	DbPrefix   string `toml:"dbPrefix"`
	AppKey     string `toml:"appKey"`
}

type wechat struct {
	CodeToSessURL string
	AppID         string
	AppSecret     string
}

type wechatPay struct {
	AppID       string
	MchID       string
	Key         string
	CallbackApi string
}

type wechatOpen struct {
	AppID       string
	AppSecret   string
	RedirectURI string
	Scope       string
}

type oss struct {
	Domain string
}

// initConfig initializes the app configuration by first setting defaults,
// then overriding settings from the app config file, then overriding
// It returns an error if any.
func InitConfig(configFile string) error {
	if configFile == "" {
		panic("config path is empty")
	}
	var err error

	if _, err = os.Stat(configFile); err != nil {
		return errors.New("config file err:" + err.Error())
	} else {
		configBytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			return errors.New("config load err:" + err.Error())
		}
		_, err = toml.Decode(string(configBytes), &Conf)
		if err != nil {
			return errors.New("config decode err:" + err.Error())
		}
	}
	return nil
}

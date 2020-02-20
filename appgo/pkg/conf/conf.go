package conf

import (
	"github.com/spf13/viper"
)

var (
	Conf config // holds the global app config.
)

type config struct {
	// 应用配置
	App   app
	Oss   oss
	Image image
}

type image struct {
	Domain string
	Path   string
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

func Init() error {
	err := viper.Unmarshal(&Conf)
	return err
}

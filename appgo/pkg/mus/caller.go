package mus

import (
	"github.com/gin-gonic/gin"
	"github.com/i2eco/egoshop/appgo/pkg/conf"
	"github.com/i2eco/egoshop/appgo/pkg/opensdk/wechatauth"
	"github.com/i2eco/muses/pkg/cache/mixcache"
	"github.com/i2eco/muses/pkg/database/mysql"
	"github.com/i2eco/muses/pkg/logger"
	"github.com/i2eco/muses/pkg/oss"
	musgin "github.com/i2eco/muses/pkg/server/gin"
	"github.com/i2eco/muses/pkg/session/ginsession"
	"github.com/jinzhu/gorm"
	"github.com/milkbobo/gopay/client"
)

var (
	Cfg        musgin.Cfg
	Logger     *logger.Client
	Gin        *gin.Engine
	Db         *gorm.DB
	WechatAuth *wechatauth.WxConfig
	Session    gin.HandlerFunc
	Oss        *oss.Client
	Mixcache   *mixcache.Client
)

// Init 初始化muses相关容器
func Init() error {
	Cfg = musgin.Config()
	Db = mysql.Caller("egoshop")
	Logger = logger.Caller("system")
	Gin = musgin.Caller()
	Oss = oss.Caller("egoshop")
	Mixcache = mixcache.Caller("egoshop")
	Session = ginsession.Caller()
	WechatAuth = &wechatauth.WxConfig{
		AppID:  conf.Conf.App.Wechat.AppID,
		Secret: conf.Conf.App.Wechat.AppSecret,
	}
	// todo 这个包代码写的好蠢。后期fork后更改，用于gopay.Pay(charge)
	client.InitWxMiniProgramClient(&client.WechatMiniProgramClient{
		AppID:      conf.Conf.App.WechatPay.AppID,
		MchID:      conf.Conf.App.WechatPay.MchID,
		Key:        conf.Conf.App.WechatPay.Key,
		PrivateKey: nil,
		PublicKey:  nil,
	})

	// todo 这个包代码写的好蠢。后期fork后更改，用于gopay.WeChatAppCallback(c.Writer, c.Request)
	client.InitWxAppClient(&client.WechatAppClient{
		AppID:      conf.Conf.App.WechatPay.AppID,
		MchID:      conf.Conf.App.WechatPay.MchID,
		Key:        conf.Conf.App.WechatPay.Key,
		PrivateKey: nil,
		PublicKey:  nil,
	})
	return nil

}

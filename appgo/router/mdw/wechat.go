package mdw

import (
	"github.com/gin-gonic/gin"

	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/pkg/mus"
	"github.com/goecology/egoshop/appgo/token"
)

var DefaultWechatUid = "github.com/goecology/egoshop/wechatuid"
var DefaultWechatUsername = "github.com/goecology/egoshop/wechatUsername"

// WechatAccess 微信登录校验中间件
func WechatAccessMustLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Access-Token")
		if !token.GetAccessor().CheckAccessToken(c, accessToken) {
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "auth token error",
			})
			c.Abort()
			return
		}

		userInfo, err := token.GetAccessor().DecodeAccessToken(accessToken)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "auth token decode error",
			})
			c.Abort()
			return
		}
		uid, flag := userInfo["sub"].(float64)
		if !flag {
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "auth token assert error",
			})
			c.Abort()
			return
		}
		c.Set(DefaultWechatUid, int(uid))
		c.Next()
	}
}

// WechatAccess 微信登录校验中间件
func WechatAccessLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Access-Token")
		if !token.GetAccessor().CheckAccessToken(c, accessToken) {
			c.Next()
			return
		}

		userInfo, err := token.GetAccessor().DecodeAccessToken(accessToken)
		if err != nil {
			c.Next()
			return
		}
		uid, flag := userInfo["sub"].(float64)
		if !flag {
			c.Next()
			return
		}
		c.Set(DefaultWechatUid, int(uid))
		c.Next()
	}
}

// WechatUid 获取微信id
func WechatUid(c *gin.Context) int {
	return c.MustGet(DefaultWechatUid).(int)
}

func WechatMaybeUid(c *gin.Context) (uid int, flag bool) {
	var uidInterface interface{}
	uidInterface, flag = c.Get(DefaultWechatUid)
	if !flag {
		return
	}
	uid = uidInterface.(int)
	return
}

// WechatUserName 获取微信昵称
func WechatUserName(c *gin.Context) (username string) {
	value, flag := c.Get(DefaultWechatUsername)
	if flag {
		username = value.(string)
		return
	}

	uid := c.MustGet(DefaultWechatUid).(int)
	user := mysql.User{}
	err := mus.Db.Where("id = ?", uid).Find(&user)
	if err != nil {
		// todo log
		return
	}
	c.Set(DefaultWechatUsername, user.Name)
	username = user.Name
	return
}

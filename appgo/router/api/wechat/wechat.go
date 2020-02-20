package wechat

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/pkg/base"
	"github.com/goecology/egoshop/appgo/pkg/code"
	"github.com/goecology/egoshop/appgo/pkg/mus"
)

type WechatUser struct {
	NickName string `json:"nickName"`
	Avatar   string `json:"avatar"`
	Intro    string `json:"intro"`
}

func Login(c *gin.Context) {
	wxCode := c.Request.Header.Get("X-WX-Code")
	wxData := c.Request.Header.Get("X-WX-Encrypted-Data")
	wxIv := c.Request.Header.Get("X-WX-IV")
	var user mysql.User
	var nickname string
	var avatar string
	if wxIv != "" {
		// 1.登录信
		userInfo, err := mus.WechatAuth.WexLogin(wxCode, wxData, wxIv)
		if err != nil {
			base.JSONErr(c, code.LoginWechatErr, err)
			return
		}
		nickname = userInfo.NickName
		avatar = userInfo.AvatarURL
		// 写user_open表
		createUserOpen := &mysql.UserOpen{
			Genre:        0,
			Name:         "",
			WechatOpenid: "",
			AppOpenid:    "",
			MiniOpenid:   userInfo.OpenID,
			MiniOpenid2:  "",
			Unionid:      userInfo.UnionID,
			AccessToken:  "",
			ExpiresIn:    0,
			RefreshToken: "",
			Scope:        "",
			Nickname:     userInfo.NickName,
			Avatar:       userInfo.AvatarURL,
			Sex:          userInfo.Gender,
			Country:      userInfo.Country,
			Province:     userInfo.Province,
			City:         userInfo.City,
			// InfoAggregate: createUserInfo,
			State:     0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		tx := mus.Db.Begin()
		user, err = dao.WechatUser.Login(c, tx, createUserOpen, "mini")
		if err != nil {
			tx.Rollback()
			base.JSONErr(c, code.MsgErr, err)
			return
		}
		tx.Commit()
	} else {
		userInfo2, err := mus.WechatAuth.LoginCode(wxCode)
		if err != nil {
			base.JSONErr(c, code.LoginWechatErr, err)
			return
		}
		nickname = userInfo2.NickName
		avatar = userInfo2.HeadImgURL
		// 写user_open表
		createUserOpen := &mysql.UserOpen{
			Genre:        0,
			Name:         "",
			WechatOpenid: "",
			AppOpenid:    "",
			MiniOpenid:   userInfo2.OpenID,
			MiniOpenid2:  "",
			Unionid:      userInfo2.UnionID,
			AccessToken:  "",
			ExpiresIn:    0,
			RefreshToken: "",
			Scope:        "",
			Nickname:     userInfo2.NickName,
			Avatar:       userInfo2.HeadImgURL,
			Sex:          int(userInfo2.Sex),
			Country:      userInfo2.Country,
			Province:     userInfo2.Province,
			City:         userInfo2.City,
			// InfoAggregate: createUserInfo,
			State:     0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		tx := mus.Db.Begin()
		user, err = dao.WechatUser.Login(c, tx, createUserOpen, "mini")
		if err != nil {
			tx.Rollback()
			base.JSONErr(c, code.MsgErr, err)
			return
		}
		tx.Commit()
	}

	// 4.创建token
	respToken, err := dao.AccessToken.CreateAccessToken(c, user.Id, time.Now().Unix())
	if err != nil {
		base.JSONErr(c, code.LoginWechatErr5, err)
		return
	}

	wechatUser := WechatUser{
		NickName: nickname,
		Avatar:   avatar,
		Intro:    "todo",
	}

	base.JSONOK(c, gin.H{
		"skey":     respToken.AccessToken,
		"userinfo": wechatUser,
	})
}

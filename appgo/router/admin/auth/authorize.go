package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/router/mdw"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
	"github.com/goecology/egoshop/appgo/pkg/base"
	"github.com/goecology/egoshop/appgo/pkg/code"
	"regexp"
)

var emailRgx = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-z]{2,4}$`)

// {status: "error", type: "account", currentAuthority: "guest"}
func Login(c *gin.Context) {
	var err error
	// 如果已经登录
	respView := trans.RespOauthLogin{
		CurrentAuthority: "admin",
	}

	if mdw.AdminAuthed(c) {
		base.JSONOK(c, respView)
		return
	}

	reqView := &trans.ReqOauthLogin{}
	err = c.Bind(reqView)
	if err != nil {
		base.JSONErr(c, code.AdminLoginParamsErr, err)
		return
	}

	// 对Identity进行校验，先判断是否是邮箱，若不是邮箱则当做用户名
	var oneUser *mysql.User
	if emailRgx.MatchString(reqView.Name) {
		oneUser, err = dao.User.GetBizByPwd("", reqView.Name, reqView.Pwd, c.ClientIP())
		if err != nil {
			base.JSONErr(c, code.AdminLoginEmailErr, err)
			return
		}
	} else {
		oneUser, err = dao.User.GetBizByPwd(reqView.Name, "", reqView.Pwd, c.ClientIP())
		if err != nil {
			base.JSONErr(c, code.AdminLoginNicknameErr, err)
			return
		}
	}
	err = mdw.UpdateUser(c, *oneUser)
	if err != nil {
		base.JSONErr(c, code.AdminLoginUpdateErr, err)
		return
	}
	base.JSONOK(c, respView)
	return
}

func Logout(c *gin.Context) {
	err := mdw.Logout(c)
	if err != nil {
		base.JSONErr(c, code.AdminLogoutErr, err)
		return
	}
	base.JSONOK(c)
	return

}

func Self(c *gin.Context) {
	resp, err := dao.User.Info(c, mdw.AdminUid(c))
	if err != nil {
		base.JSONErr(c, code.AdminUserInfoErr, err)
		return
	}
	base.JSONOK(c, resp)
	return
}

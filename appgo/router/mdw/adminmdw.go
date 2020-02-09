package mdw

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"net/http"
)

const SessionDefaultKey = "goecology/mdw/session"
const ContextUser = "mdw/member"

func init() {
	gob.Register(mysql.User{})
}

// 后台取用户
func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		one := mysql.User{}
		// 从session中获取用户信息
		if user, ok := DefaultSessionUser(c); ok && user.Id > 0 {
			one = user
		} else {

		}
		c.Set(ContextUser, one)
		c.Next()
	}
}

// 后台取用户
func LoginAPIRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		one := mysql.User{}
		// 从session中获取用户信息
		if user, ok := DefaultSessionUser(c); ok && user.Id > 0 {
			one = user
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 401,
				"data": "",
				"msg":  "user not login",
			})
			c.Abort()
			return
		}
		c.Set(ContextUser, one)
		c.Next()
	}
}

// 后台取用户
func DefaultSessionUser(c *gin.Context) (mysql.User, bool) {
	resp, flag := sessions.Default(c).Get(SessionDefaultKey).(mysql.User)
	return resp, flag
}

// 后台取用户
func DefaultContextUser(c *gin.Context) mysql.User {
	var resp mysql.User
	respI, flag := c.Get(ContextUser)
	if flag {
		resp = respI.(mysql.User)
	}
	return resp

}

// Authed 鉴权通过
func AdminAuthed(c *gin.Context) bool {
	if user, ok := DefaultSessionUser(c); ok && user.Id > 0 {
		return true
	}
	return false
}

// 后台 Uid 返回uid
func AdminUid(c *gin.Context) int {
	return DefaultContextUser(c).Id
}

// UpdateUser updates the User object stored in the session. This is useful incase a change
// is made to the user model that needs to persist across requests.
func UpdateUser(c *gin.Context, a mysql.User) error {
	s := sessions.Default(c)
	s.Options(sessions.Options{
		Path:     "/",
		MaxAge:   24 * 3600,
		Secure:   false,
		HttpOnly: true,
	})
	s.Set(SessionDefaultKey, a)
	return s.Save()
}

// Logout will clear out the session and call the Logout() user function.
func Logout(c *gin.Context) error {
	s := sessions.Default(c)
	s.Options(sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
	})

	s.Delete(SessionDefaultKey)
	return s.Save()
}

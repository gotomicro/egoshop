package dao

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/i2eco/egoshop/appgo/model/mysql"
	"github.com/i2eco/egoshop/appgo/pkg/random"
	"github.com/jinzhu/gorm"
)

const (
	GenreWechatMini = 3
)

type wechatUser struct{}

func InitWechatUser() *wechatUser {
	return &wechatUser{}
}

func (*wechatUser) Login(c *gin.Context, tx *gorm.DB, userOpenInfo *mysql.UserOpen, typeWechat string) (userInfo mysql.User, err error) {
	var respMysql mysql.UserOpen

	if typeWechat == "mini" {
		// 2.获取数据库信息
		respMysql, err = UserOpen.InfoX(c, mysql.Conds{"mini_openid": userOpenInfo.MiniOpenid})
		if err != nil && err != gorm.ErrRecordNotFound {
			return
		}
		// 当找不到记录，才需要重新创建name
		if err != nil && err == gorm.ErrRecordNotFound {
			userOpenInfo.Name = "wechat_mini_" + userOpenInfo.MiniOpenid + "_" + random.GetRandomString(8)
			userOpenInfo.Genre = 3
		} else {
			userOpenInfo.Name = respMysql.Name
		}
	} else if typeWechat == "app" {
		// 2.获取数据库信息
		respMysql, err = UserOpen.InfoX(c, mysql.Conds{"app_openid": userOpenInfo.AppOpenid})
		if err != nil && err != gorm.ErrRecordNotFound {
			return
		}

		// 当找不到记录，才需要重新创建name
		if err != nil && err == gorm.ErrRecordNotFound {
			userOpenInfo.Name = "wechat_app_" + userOpenInfo.AppOpenid + "_" + random.GetRandomString(8)
			userOpenInfo.Genre = 1
		} else {
			userOpenInfo.Name = respMysql.Name
		}

	} else {
		err = errors.New("not exist type wechat")
		return
	}

	// 如果不存在open id，查一下union id。
	if err == gorm.ErrRecordNotFound {
		respMysql, err = UserOpen.InfoX(c, mysql.Conds{"unionid": userOpenInfo.Unionid})
		if err != nil && err != gorm.ErrRecordNotFound {
			return
		}

		// 如果union id，也不存在，那么就创建用户
		if err == gorm.ErrRecordNotFound {
			err = UserOpen.Create(c, tx, userOpenInfo)
			if err != nil {
				return
			}

			newUser := &mysql.User{
				Name:          userOpenInfo.Name,
				Nickname:      userOpenInfo.Nickname,
				OpenId:        userOpenInfo.Id,
				OpenType:      1,
				Avatar:        userOpenInfo.Avatar,
				LastLoginIP:   c.ClientIP(),
				LastLoginTime: time.Now(),
				Status:        1,
				Role:          2,
			}

			// 写user表
			err = User.Create(c, tx, newUser)
			if err != nil {
				return
			}
			// 将uid更新到user open里
			err = UserOpen.Update(c, tx, userOpenInfo.Id, mysql.Ups{"uid": newUser.Id})
			if err != nil {
				return
			}
			// 如果存在union id，则更新对应的字段信息
		} else {
			userOpenInfo.Name = respMysql.Name
			if typeWechat == "mini" {
				err = UserOpen.Update(c, tx, respMysql.Id, mysql.Ups{"mini_openid": userOpenInfo.MiniOpenid})
				if err != nil {
					return
				}
			} else if typeWechat == "app" {
				err = UserOpen.Update(c, tx, respMysql.Id, mysql.Ups{"app_openid": userOpenInfo.AppOpenid})
				if err != nil {
					return
				}
			}
		}
	}

	err = User.UpdateX(c, tx, mysql.Conds{"name": userOpenInfo.Name}, mysql.Ups{
		"last_login_ip":   c.ClientIP(),
		"last_login_time": time.Now(),
	})
	if err != nil {
		return
	}

	err = tx.Where("name = ? AND status = ?", userOpenInfo.Name, 1).Find(&userInfo).Error
	if err != nil {
		return
	}
	return
}

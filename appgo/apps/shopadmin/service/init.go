package service

import "github.com/goecology/egoshop/appgo/dao"

func Init() {
	InitGen()
	InitOssCli()
	dao.WechatUser = dao.InitWechatUser()
}

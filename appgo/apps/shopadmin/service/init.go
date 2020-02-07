package service

import "github.com/goecology/egoshop/appgo/dao"

func Init() error {
	InitGen()
	InitOssCli()
	dao.WechatUser = dao.InitWechatUser()
	return nil
}

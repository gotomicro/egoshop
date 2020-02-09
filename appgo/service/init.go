package service

import (
	"github.com/goecology/egoshop/appgo/pkg/mus"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/spf13/viper"
)

var (
	QueueSignin *queueSignin
	QueuePoint  *queuePoint
	QueueView   *queueView
)

func Init() error {
	InitGen()
	InitOssCli(
		viper.GetString("oss.endpoint"),
		viper.GetString("oss.accessKeyID"),
		viper.GetString("oss.accessKeySecret"),
		viper.GetString("oss.bucket"),
		mus.Logger,
	)
	QueueSignin = InitQueueSignin()
	QueuePoint = InitQueuePoint()
	dao.WechatUser = dao.InitWechatUser()
	QueueView = InitQueueView()
	return nil
}

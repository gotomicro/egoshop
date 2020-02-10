package service

import (
	"github.com/goecology/egoshop/appgo/dao"
)

var (
	QueueSignin *queueSignin
	QueuePoint  *queuePoint
	QueueView   *queueView
)

func Init() error {
	InitGen()
	QueueSignin = InitQueueSignin()
	QueuePoint = InitQueuePoint()
	dao.WechatUser = dao.InitWechatUser()
	QueueView = InitQueueView()
	return nil
}

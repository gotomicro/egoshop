package service

import (
	"github.com/goecology/egoshop/appgo/apps/shopapi/pkg/mus"
	"github.com/goecology/egoshop/appgo/dao"
)

func InitGen() {
	dao.Address = dao.InitAddress(mus.Logger, mus.Db)
	dao.Order = dao.InitOrder(mus.Logger, mus.Db)
	dao.AddressType = dao.InitAddressType(mus.Logger, mus.Db)
	dao.ComSku = dao.InitComSku(mus.Logger, mus.Db)
	dao.ComSpec = dao.InitComSpec(mus.Logger, mus.Db)
	dao.Image = dao.InitImage(mus.Logger, mus.Db)
	dao.Cart = dao.InitCart(mus.Logger, mus.Db)
	dao.ComRelateCate = dao.InitComRelateCate(mus.Logger, mus.Db)
	dao.AccessToken = dao.InitAccessToken(mus.Logger, mus.Db)
	dao.AdminUser = dao.InitAdminUser(mus.Logger, mus.Db)
	dao.OrderGoods = dao.InitOrderGoods(mus.Logger, mus.Db)
	dao.OrderPay = dao.InitOrderPay(mus.Logger, mus.Db)
	dao.Comment = dao.InitComment(mus.Logger, mus.Db)
	dao.PointLog = dao.InitPointLog(mus.Logger, mus.Db)
	dao.UserOpen = dao.InitUserOpen(mus.Logger, mus.Db)
	dao.PointLimit = dao.InitPointLimit(mus.Logger, mus.Db)
	dao.UserGoods = dao.InitUserGoods(mus.Logger, mus.Db)
	dao.Signin = dao.InitSignin(mus.Logger, mus.Db)
	dao.User = dao.InitUser(mus.Logger, mus.Db)
	dao.Attachment = dao.InitAttachment(mus.Logger, mus.Db)
	dao.SigninLog = dao.InitSigninLog(mus.Logger, mus.Db)
	dao.ComStore = dao.InitComStore(mus.Logger, mus.Db)
	dao.ComSpecValue = dao.InitComSpecValue(mus.Logger, mus.Db)
	dao.ComImage = dao.InitComImage(mus.Logger, mus.Db)
	dao.OrderLog = dao.InitOrderLog(mus.Logger, mus.Db)
	dao.Banner = dao.InitBanner(mus.Logger, mus.Db)
	dao.Freight = dao.InitFreight(mus.Logger, mus.Db)
	dao.OrderExtend = dao.InitOrderExtend(mus.Logger, mus.Db)
	dao.ComCate = dao.InitComCate(mus.Logger, mus.Db)
	dao.Com = dao.InitCom(mus.Logger, mus.Db)

}

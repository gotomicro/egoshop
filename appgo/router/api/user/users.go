package user

import (
	"github.com/goecology/egoshop/appgo/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/pkg/mus"
	"github.com/goecology/egoshop/appgo/router/mdw"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/constx"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
	"github.com/goecology/egoshop/appgo/pkg/base"
	"github.com/goecology/egoshop/appgo/pkg/code"
	"github.com/goecology/egoshop/appgo/pkg/imagex"
	"github.com/spf13/cast"
)

// 我的基本信息
func Stats(c *gin.Context) {
	//todo 注意安全问题
	user, _ := dao.User.Info(c, mdw.WechatUid(c))
	signinData, err := dao.Signin.InfoX(c, mysql.Conds{
		"uid": mdw.WechatUid(c),
	})

	// 获取今天0点时间
	todayZeroStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", todayZeroStr)
	todayZero := t.Unix()
	// 今天是否签到
	isSignin := false

	// 连续签到天数
	signinDur := 0
	if err == nil {
		// 上次更新时间
		updatedUnix := signinData.UpdatedAt.Unix()
		// 说明是连续签到的，可以返回连续签到时间
		if (todayZero-24*60*60) < updatedUnix && todayZero >= updatedUnix {
			signinDur = signinData.SigninCnt
		}

		if updatedUnix > todayZero {
			isSignin = true
		}
	}
	totalCntCollect := 0
	totalCntRead := 0
	totalCntShare := 0
	totalCntOrderNew := 0
	totalCntOrderPay := 0
	totalCntOrderSend := 0
	uid := mdw.WechatUid(c)
	mus.Db.Table("user_goods").Where("uid = ? and is_collect = ?", mdw.WechatUid(c), 1).Count(&totalCntCollect)
	mus.Db.Table("user_goods").Where("uid = ? and is_read = ?", mdw.WechatUid(c), 1).Count(&totalCntRead)
	mus.Db.Table("user_goods").Where("uid = ? and is_share = ?", mdw.WechatUid(c), 1).Count(&totalCntShare)
	// 待付款
	mus.Db.Table("order").Where("uid = ? and state = ?", uid, constx.OrderStateNew).Count(&totalCntOrderNew)
	// 待发货
	mus.Db.Table("order").Where("uid = ? and state = ?", uid, constx.OrderStatePay).Count(&totalCntOrderPay)
	// 待收货
	mus.Db.Table("order").Where("uid = ? and state = ?", uid, constx.OrderStateSend).Count(&totalCntOrderSend)

	base.JSONOK(c, gin.H{
		"point":        user.Point,
		"signinDur":    signinDur,
		"roleName":     constx.RoleNameMap[user.Role],
		"isSignin":     isSignin,
		"cntCollect":   totalCntCollect,
		"cntRead":      totalCntRead,
		"cntShare":     totalCntShare,
		"cntOrderNew":  totalCntOrderNew,
		"cntOrderPay":  totalCntOrderPay,
		"cntOrderSend": totalCntOrderSend,
	})
}

// StarGoods 关注商品
func StarGoods(c *gin.Context) {
	goodsInfo, err := mdw.ValidateGoods(c)
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	uid := mdw.WechatUid(c)
	tx := mus.Db.Begin()
	err = dao.UserGoods.CreateOrUpdate(tx, uid, goodsInfo, "star")
	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	tx.Commit()
	service.QueuePoint.Push(service.PointTask{
		TypeData: constx.PointStar,
		Uid:      uid,
	})
	base.JSONOK(c)
}

// UnstarGoods 取消关注商品
func UnstarGoods(c *gin.Context) {
	goodsInfo, err := mdw.ValidateGoods(c)
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	tx := mus.Db.Begin()

	err = dao.UserGoods.Cancel(tx, mdw.WechatUid(c), goodsInfo, "star")
	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	tx.Commit()
	base.JSONOK(c)
}

// StarGoods 关注商品
func CollectGoods(c *gin.Context) {
	goodsInfo, err := mdw.ValidateGoods(c)
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	uid := mdw.WechatUid(c)
	tx := mus.Db.Begin()
	err = dao.UserGoods.CreateOrUpdate(tx, uid, goodsInfo, "collect")
	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	tx.Commit()
	service.QueuePoint.Push(service.PointTask{
		TypeData: constx.PointStar,
		Uid:      uid,
	})
	base.JSONOK(c)
}

// UnstarGoods 取消关注商品
func UncollectGoods(c *gin.Context) {
	goodsInfo, err := mdw.ValidateGoods(c)
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	tx := mus.Db.Begin()

	err = dao.UserGoods.Cancel(tx, mdw.WechatUid(c), goodsInfo, "collect")
	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	tx.Commit()
	base.JSONOK(c)
}

// ShareGoods 分享商品
func ShareGoods(c *gin.Context) {
	goodsInfo, err := mdw.ValidateGoods(c)
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	uid := mdw.WechatUid(c)

	tx := mus.Db.Begin()
	err = dao.UserGoods.CreateOrUpdate(tx, uid, goodsInfo, "share")
	if err != nil {
		tx.Rollback()
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	tx.Commit()
	service.QueuePoint.Push(service.PointTask{
		TypeData: constx.PointStar,
		Uid:      uid,
	})
	base.JSONOK(c)
}

// GetOne 获取指定用户信息
func GetOne(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	user, _ := dao.User.Info(c, id)
	base.JSONOK(c, user)
}

// todo 迁移到model
type OutputUserGoods struct {
	Id          int       `json:"id" form:"id" gorm:"primary_key"` // 主键ID
	CreatedAt   time.Time `json:"createdAt" form:"createdAt" `     // 创建时间
	GoodsId     int       `json:"goodsId" form:"goodsId" `         // 商品ID
	TypeId      int       `json:"typeId" form:"typeId" `           // 类型id
	Name        string    `json:"name" form:"name" `               // 商品名称
	Cover       string    `json:"cover"`
	IsPay       int       `json:"isPay" form:"isPay" `       // 是否购买
	IsStar      int       `json:"isStar" form:"isStar" `     // 是否关注
	IsCreate    int       `json:"isCreate" form:"isCreate" ` // 是否上传
	IsShare     int       `json:"isShare" form:"isShare" `   // 是否分享
	IsPrePay    int       `json:"isPrePay" form:"isPrePay" ` // 是否预购买
	CntView     int       `json:"cntView"`
	Desc        string    `json:"desc"`
	CreatedName string    `json:"createdName"`
	TypeName    string    `json:"typeName"`
}

// 查看用户自己购买，收藏等商品
func List(c *gin.Context) {

	req := ReqList{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request error")
		return
	}

	reqPage := trans.ReqPage{
		Current:  req.Current,
		PageSize: req.PageSize,
		Sort:     "",
	}

	typeName, flag := constx.GoodsTypeMap[req.Tid]
	if !flag {
		base.JSON(c, code.MsgErr)
		return
	}

	conds := mysql.Conds{
		"uid":    mdw.WechatUid(c),
		typeName: 1,
	}
	if req.TypeId != 0 {
		conds["type_id"] = req.TypeId
	}

	total, resp := dao.UserGoods.ListPage(c, conds, &reqPage)

	output := make([]OutputUserGoods, 0)

	//// todo 同一类型，用in语句，合并io操作
	for _, value := range resp {
		switch value.TypeId {
		case constx.TypeCom:
			info, _ := dao.Com.Info(c, value.GoodsId)
			userInfo, _ := dao.User.Info(c, info.CreatedBy)
			output = append(output, OutputUserGoods{
				Id:          value.Id,
				CreatedAt:   value.CreatedAt,
				GoodsId:     value.GoodsId,
				TypeId:      value.TypeId,
				Name:        info.Title,
				Cover:       imagex.ShowImg(info.Cover, "x1"),
				CreatedName: userInfo.Nickname,
				TypeName:    "商品",
			})
		}
	}

	base.JSONWechatList(c, output, total, reqPage.PageSize)
}

func GetUserGoodsOne(c *gin.Context) {
	gid := cast.ToInt(c.Param("gid"))
	typeId := cast.ToInt(c.Param("typeId"))
	userGoods, _ := dao.UserGoods.InfoX(c, mysql.Conds{
		"uid":      mdw.WechatUid(c),
		"goods_id": gid,
		"type_id":  typeId,
	})
	base.JSONOK(c, userGoods)
}

type File struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Key  string `json:"key"`
}

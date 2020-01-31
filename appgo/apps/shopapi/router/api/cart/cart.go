package cart

import (
	"errors"
	"time"

	"github.com/goecology/egoshop/appgo/model/trans"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/apps/shopapi/pkg/mus"
	"github.com/goecology/egoshop/appgo/apps/shopapi/router/mdw"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/pkg/base"
	"github.com/goecology/egoshop/appgo/pkg/code"
	"github.com/jinzhu/gorm"
)

func Info(c *gin.Context) {
	req := ReqInfo{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	resp, err := dao.Cart.InfoX(c, mysql.Conds{
		"created_by": mdw.WechatUid(c),
		"com_sku_id": req.ComSkuId,
	})
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	base.JSONOK(c, resp)

}

func List(c *gin.Context) {
	req := ReqList{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	reqPage := trans.ReqPage{
		Current:  req.Current,
		PageSize: 10,
		Sort:     "updated_at desc",
	}

	uid := mdw.WechatUid(c)

	list, total := dao.Cart.ListAddition(c, uid, req.CartIds, &reqPage)
	base.JSONWechatList(c, list, total, reqPage.PageSize)

}

func Create(c *gin.Context) {
	req := ReqCreate{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	var (
		cartInfo mysql.Cart
		err      error
		codeInt  int
	)
	uid := mdw.WechatUid(c)

	// 添加商品
	if req.TypeId == 1 {
		cartInfo, err, codeInt = addCom(c, req)
		if err != nil {
			base.JSONErr(c, codeInt, err)
			return
		}
	} else {
		base.JSONErr(c, code.MsgErr, nil)
		return
	}
	var cnt int
	mus.Db.Model(mysql.Cart{}).Where("created_by = ?", uid).Count(&cnt)
	base.JSON(c, code.MsgOk, gin.H{
		"cnt":  cnt,
		"info": cartInfo,
	})
}

func addCom(c *gin.Context, req ReqCreate) (cartInfo mysql.Cart, err error, codeInt int) {
	codeInt = 1
	uid := mdw.WechatUid(c)
	cartInfo, _ = dao.Cart.InfoX(c, mysql.Conds{
		"com_sku_id": req.ComSkuId,
		"created_by": uid,
	})

	comSkuInfo, err := dao.ComSku.InfoX(c, mysql.Conds{
		"id": req.ComSkuId,
	})
	if err != nil {
		codeInt = code.ComSkuInfoSystemError
		err = errors.New("商品信息不存在")
		return
	}

	comInfo, err := dao.Com.InfoX(c, mysql.Conds{
		"id":         comSkuInfo.ComId,
		"is_on_sale": 1,
		"sale_time":  mysql.Cond{"<", time.Now()},
	})
	if err != nil {
		codeInt = code.ComInfoSystemError
		return
	}

	// 购物车里原基础增加数量
	comNum := req.Num
	if cartInfo.Num > 0 {
		comNum += cartInfo.Num
	}

	if comInfo.Id == 0 {
		codeInt = code.ComInfoNotExist
		err = errors.New("产品不存在")
		return
	}

	if comInfo.Stock < 1 {
		err = errors.New("商品库存不足")
		return
	}

	if comInfo.Stock < comNum {
		err = errors.New("库存不足")
		return
	}

	if cartInfo.Id > 0 {
		err = dao.Cart.UpdateX(c, mus.Db, mysql.Conds{
			"com_sku_id": req.ComSkuId,
			"created_by": uid,
		}, mysql.Ups{
			"num":        gorm.Expr("num + ?", req.Num),
			"updated_by": uid,
		})
	} else {
		cartInfo = mysql.Cart{
			CreatedBy: uid,
			UpdatedBy: uid,
			ComSkuId:  req.ComSkuId,
			Num:       req.Num,
			TypeId:    2,
		}
		err = dao.Cart.Create(c, mus.Db, &cartInfo)
	}

	if err != nil {
		return
	}
	return
}

func Update(c *gin.Context) {
	req := ReqUpdate{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request app list params is error")
		return
	}

	uid := mdw.WechatUid(c)
	cartInfo, _ := dao.Cart.InfoX(c, mysql.Conds{
		"com_sku_id": req.ComSkuId,
		"created_by": uid,
	})

	if cartInfo.Id == 0 {
		base.JSON(c, code.MsgErr, "购物车不存在")
		return
	}

	comSkuInfo, _ := dao.ComSku.InfoX(c, mysql.Conds{
		"id": req.ComSkuId,
	})

	if comSkuInfo.Id == 0 {
		base.JSON(c, code.MsgErr, "产品不存在")
		return
	}

	comInfo, _ := dao.Com.InfoX(c, mysql.Conds{
		"id":         comSkuInfo.ComId,
		"is_on_sale": 1,
		"sale_time":  mysql.Cond{"<", time.Now()},
	})
	if comInfo.Id == 0 {
		base.JSON(c, code.MsgErr, "产品不存在")
		return
	}

	// 库存小于传过来的参数
	if comSkuInfo.Stock == 0 {
		base.JSON(c, code.MsgErr, "库存不足")
		return
	}

	if comSkuInfo.Stock < req.Quantity {
		err := dao.Cart.UpdateX(c, mus.Db, mysql.Conds{
			"id":         cartInfo.Id,
			"updated_by": uid,
		}, mysql.Ups{
			"num": comSkuInfo.Stock,
		})
		if err != nil {
			base.JSON(c, code.MsgErr, "更新失败")
			return
		}
		base.JSON(c, code.MsgErr, "库存较少，更新给用户")
		return
	}

	err := dao.Cart.UpdateX(c, mus.Db, mysql.Conds{
		"id":         cartInfo.Id,
		"updated_by": uid,
	}, mysql.Ups{
		"num":        req.Quantity,
		"com_sku_id": comSkuInfo.Id,
	})
	if err != nil {
		base.JSON(c, code.MsgErr, "更新失败")
		return
	}
	var cnt int
	mus.Db.Model(mysql.Cart{}).Where("created_by = ?", uid).Count(&cnt)
	base.JSON(c, code.MsgOk, gin.H{
		"cnt":  cnt,
		"info": cartInfo,
	})

}

func Del(c *gin.Context) {
	req := ReqDel{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request app list params is error")
		return
	}

	var (
		comSkuIds []int
		feedIds   []int
	)

	for _, info := range req.Ids {
		if info.TypeId == 1 {
			feedIds = append(feedIds, info.FeedId)
		} else if info.TypeId == 2 {
			comSkuIds = append(comSkuIds, info.ComSkuIds)
		}
	}

	uid := mdw.WechatUid(c)

	tx := mus.Db.Begin()

	err := dao.Cart.DeleteX(c, tx, mysql.Conds{
		"updated_by": uid,
		"com_sku_id": mysql.Cond{"in", comSkuIds},
	})

	if err != nil {
		tx.Rollback()
		base.JSON(c, code.MsgErr, "删除失败")
		return
	}

	err = dao.Cart.DeleteX(c, tx, mysql.Conds{
		"updated_by": uid,
		"feed_id":    mysql.Cond{"in", feedIds},
	})

	if err != nil {
		tx.Rollback()
		base.JSON(c, code.MsgErr, "删除失败")
		return
	}
	tx.Commit()
	base.JSON(c, code.MsgOk, "删除成功")
}
func Check(c *gin.Context) {
	req := ReqCheck{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request app list params is error")
		return
	}

	uid := mdw.WechatUid(c)
	err := dao.Cart.UpdateX(c, mus.Db, mysql.Conds{
		"updated_by": uid,
		"com_sku_id": mysql.Cond{"in", req.ComSkuIds},
	}, mysql.Ups{
		"is_check": req.IsCheck,
	})

	if err != nil {
		base.JSON(c, code.MsgErr)
		return
	}
	base.JSON(c, code.MsgOk)
}

func Exist(c *gin.Context) {
	req := ReqInfo{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request app list params is error")
		return
	}

	_, err := dao.Cart.InfoX(c, mysql.Conds{
		"created_by": mdw.WechatUid(c),
	})

	if err != nil {
		base.JSON(c, code.MsgErr)
		return
	}

	base.JSON(c, code.MsgOk)

}

func TotalNum(c *gin.Context) {
	uid := mdw.WechatUid(c)
	var cnt int
	mus.Db.Model(mysql.Cart{}).Where("created_by = ?", uid).Count(&cnt)
	base.JSON(c, code.MsgOk, cnt)
}

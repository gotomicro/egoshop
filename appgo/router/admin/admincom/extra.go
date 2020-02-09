package admincom

import (
	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/pkg/mus"
	"github.com/goecology/egoshop/appgo/router/mdw"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
	"github.com/goecology/egoshop/appgo/pkg/base"
	"github.com/goecology/egoshop/appgo/pkg/code"
)

func OnSale(c *gin.Context) {
	req := ReqOnSale{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	tx := mus.Db.Begin()
	for _, id := range req.Ids {
		err := dao.Com.Update(c, tx, id, mysql.Ups{
			"is_on_sale": 1,
			"updated_by": mdw.AdminUid(c),
		})
		if err != nil {
			tx.Rollback()
			base.JSON(c, code.MsgErr, "update error")
			return
		}
	}
	tx.Commit()
	reqPage := &trans.ReqPage{}
	total, list := dao.Com.ListPage(c, mysql.Conds{}, reqPage)
	base.JSONList(c, list, reqPage.Current, reqPage.PageSize, total)
}

func OffSale(c *gin.Context) {
	req := ReqOnSale{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	tx := mus.Db.Begin()
	for _, id := range req.Ids {
		err := dao.Com.Update(c, tx, id, mysql.Ups{
			"is_on_sale": 0,
			"updated_by": mdw.AdminUid(c),
		})
		if err != nil {
			tx.Rollback()
			base.JSON(c, code.MsgErr, "update error")
			return
		}
	}
	tx.Commit()
	reqPage := &trans.ReqPage{}
	total, list := dao.Com.ListPage(c, mysql.Conds{}, reqPage)
	base.JSONList(c, list, reqPage.Current, reqPage.PageSize, total)
}

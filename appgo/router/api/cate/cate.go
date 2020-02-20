package cate

import (
	//"github.com/goecology/egoshop/appgo/model/trans"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
	"github.com/goecology/egoshop/appgo/pkg/base"
	"github.com/goecology/egoshop/appgo/pkg/code"
	// "github.com/goecology/egoshop/appgo/router/mdw"
)

// 查找某分类下子分类（其下1级）
func List(c *gin.Context) {
	req := ReqList{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	if req.CateChild {
		result := dao.ComCate.ListChild(c, mysql.Conds{"pid": req.CateId}, &trans.ReqPage{Sort: "sort desc"})
		base.JSONOK(c, result)
	} else {
		Single(c)
	}
}

// 查找分类详情
func Single(c *gin.Context) {
	req := ReqList{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	if result, err := dao.ComCate.Info(c, req.CateId); err != nil {
		base.JSONErr(c, code.MsgErr, err)
	} else {
		base.JSONOK(c, result)
	}
}

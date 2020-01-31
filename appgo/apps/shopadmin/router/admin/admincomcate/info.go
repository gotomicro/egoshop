package admincomcate

import (
	"time"

	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/mdw"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/pkg/mus"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/pkg/base"
	"github.com/goecology/egoshop/appgo/pkg/code"
)

func List(c *gin.Context) {
	req := ReqList{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}
	conds := mysql.Conds{}

	// 搜索名称
	if req.Name != "" {
		conds["name"] = mysql.Cond{
			"like",
			req.Name,
		}
	}
	list, _ := dao.ComCate.List(c, conds, "")
	base.JSONList(c, list, 1, 10000, len(list))
}

func Info(c *gin.Context) {
	req := ReqInfo{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}
	// TODO 完善ReqGoodscategoryInfo字段
	res, err := dao.ComCate.Info(c, req.Id)
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	base.JSON(c, code.MsgOk, res)
}

func Create(c *gin.Context) {
	req := ReqCreate{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}

	err := dao.ComCate.Create(c, mus.Db, &mysql.ComCate{
		CreatedBy: mdw.AdminUid(c),
		UpdatedBy: mdw.AdminUid(c),
		Name:      req.Name,
		Pid:       req.Pid,
		Icon:      req.Icon,
	})

	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}

	list, _ := dao.ComCate.List(c, mysql.Conds{}, "")
	base.JSONList(c, list, 1, 10000, len(list))
}

func Update(c *gin.Context) {
	req := ReqUpdate{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}

	err := dao.ComCate.Update(c, mus.Db, req.Id, mysql.Ups{
		"name":       req.Name,
		"icon":       req.Icon,
		"updated_by": mdw.AdminUid(c),
	})
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	list, _ := dao.ComCate.List(c, mysql.Conds{}, "")
	base.JSONList(c, list, 1, 10000, len(list))
}

func Remove(c *gin.Context) {
	req := ReqRemove{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}

	err := dao.ComCate.Update(c, mus.Db, req.Id, mysql.Ups{
		"deleted_at": time.Now(),
		"updated_by": mdw.AdminUid(c),
	})
	if err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	list, _ := dao.ComCate.List(c, mysql.Conds{}, "")
	base.JSONList(c, list, 1, 10000, len(list))
}

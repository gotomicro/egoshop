package adminuser

import (
	"github.com/gin-gonic/gin"
	"github.com/i2eco/egoshop/appgo/dao"
	"github.com/i2eco/egoshop/appgo/model/mysql"
	"github.com/i2eco/egoshop/appgo/pkg/base"
	"github.com/i2eco/egoshop/appgo/pkg/code"
)

func List(c *gin.Context) {
	req := ReqList{}
	if err := c.Bind(&req); err != nil {
		base.JSON(c, code.MsgErr, "request is error")
		return
	}
	total, list := dao.User.ListPage(c, mysql.Conds{}, &req.ReqPage)
	base.JSONList(c, list, req.ReqPage.Current, req.ReqPage.PageSize, total)
}

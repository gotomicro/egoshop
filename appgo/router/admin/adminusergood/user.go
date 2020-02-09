package adminusergood

import (
	"github.com/gin-gonic/gin"
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
	total, list := dao.UserGoods.ListPage(c, mysql.Conds{}, &req.ReqPage)
	base.JSONList(c, list, req.ReqPage.Current, req.ReqPage.PageSize, total)
}

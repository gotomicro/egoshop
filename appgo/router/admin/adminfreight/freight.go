package adminfreight

import (
	"github.com/gin-gonic/gin"
	"github.com/i2eco/egoshop/appgo/dao"
	"github.com/i2eco/egoshop/appgo/model/mysql"
	"github.com/i2eco/egoshop/appgo/model/trans"
	"github.com/i2eco/egoshop/appgo/pkg/base"
)

func List(c *gin.Context) {
	total, list := dao.Freight.ListPage(c, mysql.Conds{}, &trans.ReqPage{Sort: "update_time desc"})
	base.JSONList(c, list, 0, 0, total)
}

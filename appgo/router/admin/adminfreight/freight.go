package adminfreight

import (
	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
	"github.com/goecology/egoshop/appgo/pkg/base"
)

func List(c *gin.Context) {
	total, list := dao.Freight.ListPage(c, mysql.Conds{}, &trans.ReqPage{Sort: "update_time desc"})
	base.JSONList(c, list, 0, 0, total)
}

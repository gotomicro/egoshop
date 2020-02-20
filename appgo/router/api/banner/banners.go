package banner

import (
	"github.com/gin-gonic/gin"
	"github.com/i2eco/egoshop/appgo/dao"
	"github.com/i2eco/egoshop/appgo/model/mysql"
)

func Get(c *gin.Context) {
	dao.Banner.List(c, mysql.Conds{})

	//base.JSONList(c, banners, len(banners))
}

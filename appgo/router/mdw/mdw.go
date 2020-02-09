package mdw

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/model/common"
	"github.com/goecology/egoshop/appgo/model/constx"
)

// 验证商品
func ValidateGoods(c *gin.Context) (g common.Goods, err error) {
	err = c.Bind(&g)
	if err != nil {
		return
	}
	return ValidateGoodsByParam(g)
}

func ValidateGoodsByParam(gp common.Goods) (g common.Goods, err error) {
	if gp.Gid <= 0 {
		err = errors.New("id is error")
		return gp, err
	}

	// 如果不为feed流类型和资料类型，直接报错
	if gp.TypeId != constx.TypeCom {
		err = errors.New("type is error")
		return gp, err
	}
	return gp, err
}

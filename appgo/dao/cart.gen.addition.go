package dao

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
	"github.com/goecology/egoshop/appgo/pkg/imagex"
)

func (g *cart) ListAddition(c *gin.Context, uid int, ids []int, reqList *trans.ReqPage) (output []mysql.Cart, cnt int) {
	condition := mysql.Conds{
		"created_by": uid,
	}

	if len(ids) > 0 {
		condition["id"] = mysql.Cond{"in", ids}
	}

	cnt, cartList := g.ListPage(c, condition, reqList)
	var comSkuIds []int
	for _, value := range cartList {
		if value.TypeId == 1 {
			comSkuIds = append(comSkuIds, value.ComSkuId)
		}
	}

	var (
		skuMap     map[int]mysql.ComSku
		comMap     map[int]mysql.Com
		freightMap map[int]mysql.Freight
	)

	if len(comSkuIds) > 0 {
		skuMap, _ = ComSku.ListMap(c, mysql.Conds{
			"id": mysql.Cond{"in", comSkuIds},
		})

		var comIds []int
		for _, value := range skuMap {
			comIds = append(comIds, value.ComId)
		}

		comMap, _ = Com.ListMap(c, mysql.Conds{
			"id": mysql.Cond{"in", comIds},
		})

		var freightIds []int
		for _, value := range comMap {
			freightIds = append(freightIds, value.FreightId)
		}
		freightMap, _ = Freight.ListMap(c, mysql.Conds{
			"id": mysql.Cond{"in", freightIds},
		})
	}

	output = make([]mysql.Cart, 0)
	for _, cartInfo := range cartList {
		if cartInfo.TypeId == 1 {
			skuInfo := skuMap[cartInfo.ComSkuId]
			comInfo := comMap[skuInfo.ComId]
			freightInfo := freightMap[comInfo.FreightId]

			cartInfo.PriceType = 1
			cartInfo.ComId = skuInfo.ComId
			cartInfo.Title = comInfo.Title
			cartInfo.SubTitle = comInfo.SubTitle
			cartInfo.ComIsOnSale = comInfo.IsOnSale
			cartInfo.ComFreightFee = comInfo.FreightFee
			cartInfo.ComFreightId = comInfo.FreightId
			cartInfo.PayType = comInfo.PayType
			cartInfo.Price = skuInfo.Price
			cartInfo.ComSpec = skuInfo.Spec
			cartInfo.ComWeight = skuInfo.Weight
			cartInfo.Stock = skuInfo.Stock
			cartInfo.Cover = imagex.ShowImg(skuInfo.Cover, "x1")
			cartInfo.ComFreightAreas = freightInfo.Areas

			// todo 统一所有的商品的售卖
			if comInfo.IsOnSale == 1 && time.Now().Sub(comInfo.SaleTime).Seconds() > 0 {
				// 购物车是生效状态，可以购买
				cartInfo.Status = 1
			} else {
				// 购物车是失效状态，无法购买
				cartInfo.Status = 2
			}
		}

		output = append(output, cartInfo)
	}
	return
}

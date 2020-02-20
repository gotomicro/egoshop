package apicom

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/dao"
	"github.com/goecology/egoshop/appgo/model/constx"
	"github.com/goecology/egoshop/appgo/model/mysql"
	"github.com/goecology/egoshop/appgo/model/trans"
	"github.com/goecology/egoshop/appgo/pkg/base"
	"github.com/goecology/egoshop/appgo/pkg/code"
	"github.com/goecology/egoshop/appgo/pkg/mus"
	"github.com/goecology/egoshop/appgo/router/mdw"
	"github.com/goecology/egoshop/appgo/service"
	"github.com/goecology/muses/pkg/cache/mixcache"
	"github.com/spf13/cast"
)

// 商城首页
func Index(c *gin.Context) {

	// 查询banners
	banners, _ := dao.Banner.List(c, mysql.Conds{})

	for idx, value := range banners {
		value.Image = mus.Oss.ShowImg(value.Image, "")
		banners[idx] = value
	}

	cates, _ := dao.ComCate.List(c, mysql.Conds{"status": 1})
	newProduct, err := mixcache.String(mus.Mixcache.Get("home:new"))
	if err != nil {
		newProduct = "[]"
	}
	base.JSONOK(c, gin.H{
		"banners":    banners,
		"cates":      cates,
		"newProduct": json.RawMessage(newProduct),
	})
}

func List(c *gin.Context) {
	req := ReqList{}
	if err := c.Bind(&req); err != nil {
		base.JSONErr(c, code.MsgErr, err)
		return
	}
	var reqPage trans.ReqPage
	if req.IsNew == 1 {
		reqPage = trans.ReqPage{
			Current:  req.CurrentPage,
			PageSize: 10,
			Sort:     "created_at desc",
		}
	} else if req.Order == "desc" {
		reqPage = trans.ReqPage{
			Current:  req.CurrentPage,
			PageSize: 10,
			Sort:     "price desc",
		}
	} else if req.Order == "asc" {
		reqPage = trans.ReqPage{
			Current:  req.CurrentPage,
			PageSize: 10,
			Sort:     "price asc",
		}
	} else {
		reqPage = trans.ReqPage{
			Current:  req.CurrentPage,
			PageSize: 10,
			Sort:     "sale_num desc",
		}
	}

	total, resp := dao.Com.ListPage(c, mysql.Conds{
		"is_on_sale": 1,
		"sale_time": mysql.Cond{
			"<",
			time.Now(),
		},
	}, &reqPage)
	for idx, value := range resp {
		value.Cover = mus.Oss.ShowImg(value.Cover, "")
		resp[idx] = value
	}
	base.JSONWechatList(c, resp, total, reqPage.PageSize)
}

func Info(c *gin.Context) {
	comInfo, err := dao.Com.Info(c, cast.ToInt(c.Param("id")))
	if err != nil {
		base.JSON(c, code.MsgErr)
		return
	}
	comInfo.Gallery = mus.Oss.ShowImgArr(comInfo.Gallery, "")
	comInfo.Cover = mus.Oss.ShowImg(comInfo.Cover, "")
	uid, flag := mdw.WechatMaybeUid(c)
	// 如果存在用户登录，才记录浏览记录
	if flag {
		service.QueueView.Push(service.ViewTask{
			GoodsId: comInfo.Id,
			Uid:     uid,
			TypeId:  constx.TypeCom,
			Name:    comInfo.Title,
		})
	}

	skuList, _ := dao.ComSku.List(c, mysql.Conds{
		"com_id": comInfo.Id,
	}, "id desc")
	comInfo.SkuList = skuList
	base.JSON(c, code.MsgOk, gin.H{
		"info": comInfo,
	})
}

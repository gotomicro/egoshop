package router

import (
	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/apps/shopapi/pkg/bootstrap"
	"github.com/goecology/egoshop/appgo/apps/shopapi/pkg/mus"
	"github.com/goecology/egoshop/appgo/apps/shopapi/router/api/address"
	"github.com/goecology/egoshop/appgo/apps/shopapi/router/api/apicom"
	"github.com/goecology/egoshop/appgo/apps/shopapi/router/api/buy"
	"github.com/goecology/egoshop/appgo/apps/shopapi/router/api/cart"
	"github.com/goecology/egoshop/appgo/apps/shopapi/router/api/comment"
	"github.com/goecology/egoshop/appgo/apps/shopapi/router/api/order"
	"github.com/goecology/egoshop/appgo/apps/shopapi/router/api/signin"
	apiuser "github.com/goecology/egoshop/appgo/apps/shopapi/router/api/user"
	"github.com/goecology/egoshop/appgo/apps/shopapi/router/api/wechat"
	"github.com/goecology/egoshop/appgo/apps/shopapi/router/mdw"
)

func InitRouter() *gin.Engine {
	r := mus.Gin
	apiGrp(r) // 小程序api路由组
	return r
}

func apiGrp(r *gin.Engine) {
	apiGrp := r.Group("/api/v1")
	apiGrp.POST("/wechat/callbackdev/payed_notify", buy.WXPayedNotify) // 微信支付回调
	apiGrp.POST("/wechat/callback/payed_notify", buy.WXPayedNotify)    // 微信支付回调
	apiGrp.GET("/auth/login", wechat.Login)                            // 微信登录

	noAuthComGrp := apiGrp.Group("/com")
	noAuthComGrp.Use(mdw.WechatAccessLogin())
	{
		// feed列表
		noAuthComGrp.GET("/index", apicom.Index)
		// feed列表
		noAuthComGrp.GET("/list", apicom.List)
		// feed单个信息列表
		noAuthComGrp.GET("/info/:id", apicom.Info)
	}

	// 评论系统
	noAuthCommentGrp := apiGrp.Group("/comment")
	{
		noAuthCommentGrp.GET("/list", comment.List)
		noAuthCommentGrp.GET("/listTop3", comment.ListTop3)
	}

	// 登录接口，不能加载这个中间件。否则死循环，登录不了
	if !bootstrap.Arg.Local {
		apiGrp.Use(mdw.WechatAccessMustLogin())
	}

	usersGrp := apiGrp.Group("/users")
	{
		// gocn 功能
		usersGrp.GET("/goods/list", apiuser.List) // 获取用户"关注"或"上传"所有商品
		// 签到功能
		usersGrp.POST("/signin", signin.Index)
		usersGrp.POST("/star", apiuser.StarGoods)           // 关注商品
		usersGrp.POST("/unstar", apiuser.UnstarGoods)       // 取消关注商品
		usersGrp.POST("/collect", apiuser.CollectGoods)     // 关注商品
		usersGrp.POST("/uncollect", apiuser.UncollectGoods) // 取消关注商品
		usersGrp.POST("/share", apiuser.ShareGoods)         // 分享商品
	}

	// 评论系统
	commentGrp := apiGrp.Group("/comment")
	{
		commentGrp.POST("/create", comment.Create)
	}

	// 支付系统
	payGrp := apiGrp.Group("/pay")
	{
		payGrp.POST("/buy", buy.Pay) // 购买商品
		//payGrp.POST("/rebuy", buy.RePay)         // 继续购买
		payGrp.POST("/calculate", buy.Calculate) // 继续购买
	}

	// 购物车
	cartGrp := apiGrp.Group("/cart")
	{
		cartGrp.POST("/list", cart.List)
		cartGrp.POST("/create", cart.Create)
		cartGrp.POST("/update", cart.Update)
		cartGrp.POST("/del", cart.Del)

		cartGrp.GET("/exist", cart.Exist)
		cartGrp.GET("/info", cart.Info)
		cartGrp.POST("/check", cart.Check)
		cartGrp.GET("/totalNum", cart.TotalNum)
	}

	addressGrp := apiGrp.Group("/address")
	{
		addressGrp.POST("/setDefault", address.SetDefault)
		addressGrp.GET("/default", address.Default)
		addressGrp.GET("/list", address.List)
		addressGrp.GET("/info", address.Info)
		addressGrp.POST("/create", address.Create)
		addressGrp.POST("/update", address.Update)
		addressGrp.POST("/delete", address.Del)
		addressGrp.GET("/typeList", address.TypeList)

	}

	orderGrp := apiGrp.Group("/order")
	{
		{
			orderGrp.GET("/stateNum", order.StateNum)
			orderGrp.POST("/create", order.Create)
			orderGrp.POST("/delete", order.Delete)
			orderGrp.POST("/pay", buy.Pay)
			orderGrp.GET("/list", order.List)
			orderGrp.GET("/info", order.Info)
			//orderGrp.POST("/cancel", order.Cancel)
			//orderGrp.POST("/confirmReceipt", order.ConfirmReceipt)
			//orderGrp.POST("/logistics", order.Logistics) // 物流查询
			//orderGrp.GET("/goodsList", order.GoodsList)
			//orderGrp.GET("/goodsInfo", order.GoodsInfo)
		}

		// 不需要强制登录获取的数据
		apiUsersGrp := r.Group("/api/v1/users")
		apiUsersGrp.Use(mdw.WechatAccessLogin())
		{
			// /api/v1/users/goods/stats/:gid/:typeId
			apiUsersGrp.GET("/goods/stats/:gid/:typeId", apiuser.GetUserGoodsOne) // 获取用户单个商品的状态
			apiUsersGrp.GET("/stats", apiuser.Stats)                              // 获取用户基本信息

		}
	}
}

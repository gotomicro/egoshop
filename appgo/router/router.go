package router

import (
	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/command"
	"github.com/goecology/egoshop/appgo/pkg/mus"
	"github.com/goecology/egoshop/appgo/router/admin/admincom"
	"github.com/goecology/egoshop/appgo/router/admin/admincomcate"
	"github.com/goecology/egoshop/appgo/router/admin/admincomment"
	"github.com/goecology/egoshop/appgo/router/admin/admincomspec"
	"github.com/goecology/egoshop/appgo/router/admin/admineditor"
	"github.com/goecology/egoshop/appgo/router/admin/adminfreight"
	"github.com/goecology/egoshop/appgo/router/admin/adminuser"
	"github.com/goecology/egoshop/appgo/router/admin/adminusergood"
	"github.com/goecology/egoshop/appgo/router/admin/auth"
	"github.com/goecology/egoshop/appgo/router/admin/image"
	"github.com/goecology/egoshop/appgo/router/api/address"
	"github.com/goecology/egoshop/appgo/router/api/apicom"
	"github.com/goecology/egoshop/appgo/router/api/buy"
	"github.com/goecology/egoshop/appgo/router/api/cart"
	"github.com/goecology/egoshop/appgo/router/api/cate"
	"github.com/goecology/egoshop/appgo/router/api/comment"
	"github.com/goecology/egoshop/appgo/router/api/order"
	"github.com/goecology/egoshop/appgo/router/api/signin"
	apiuser "github.com/goecology/egoshop/appgo/router/api/user"
	"github.com/goecology/egoshop/appgo/router/api/wechat"
	"github.com/goecology/egoshop/appgo/router/mdw"
	"github.com/spf13/viper"
)

func InitRouter() *gin.Engine {
	r := mus.Gin

	if command.Mode == "all" || command.Mode == "api" {
		apiGrp(r) // 小程序api路由组
	}
	if command.Mode == "all" || command.Mode == "admin" {
		adminGrp(r) // 后台admin路由组
	}

	r.StaticFile("/", viper.GetString("app.adminant.html"))
	// todo gin is sb
	r.NoRoute(func(c *gin.Context) {
		c.File(viper.GetString("app.adminant.html"))
	})
	r.Static("/public/", viper.GetString("app.adminant.public"))
	r.Static("/static/", viper.GetString("app.adminant.static"))
	r.Static("/"+viper.GetString("app.osspic"), viper.GetString("app.osspic"))
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

	apiGrp.Use(mdw.WechatAccessMustLogin())

	usersGrp := apiGrp.Group("/users")
	{
		usersGrp.GET("/goods/list", apiuser.List) // 获取用户"关注"或"上传"所有商品
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

		// 分类
		cateGrp := apiGrp.Group("/cate")
		{
			/* req的参数
			type ReqList struct {
				CateId    int  `json:"cateid"`
				CateChild bool `json:"catechild"` //是否获取子类
			}
			*/
			cateGrp.POST("/find", cate.List) // 购买商品
			//payGrp.POST("/rebuy", buy.RePay)         // 继续购买
		}
	}
}

func adminGrp(r *gin.Engine) {
	adGrp := r.Group("/admin", mus.Session)
	adGrp.POST("/auth/login", mus.Session, auth.Login)

	// TODO 路由规范
	// 权限模块
	authGrp := adGrp.Group("/auth")
	authGrp.Use(mdw.LoginAPIRequired())
	{
		authGrp.GET("/self", auth.Self)
		authGrp.GET("/logout", auth.Logout)
	}

	// 参考 https://github.com/gin-gonic/gin#parameters-in-path
	// TODO 路由规范

	// 用户模块
	usersGrp := adGrp.Group("/users")
	usersGrp.Use(mdw.LoginAPIRequired())
	{
		usersGrp.GET("/list", adminuser.List) // 获取所有用户
	}

	commentGrp := adGrp.Group("/comment")
	commentGrp.Use(mdw.LoginAPIRequired())
	{
		commentGrp.GET("/list", admincomment.List) // 获取所有用户
	}

	userfeedGrp := adGrp.Group("/usergood")
	userfeedGrp.Use(mdw.LoginAPIRequired())
	{
		userfeedGrp.GET("/list", adminusergood.List) // 获取所有用户
	}

	comGrp := adGrp.Group("/com")
	comGrp.Use(mdw.LoginAPIRequired())
	{
		comGrp.GET("/list", admincom.List)
		comGrp.GET("/one/:id", admincom.One)
		comGrp.POST("/create", admincom.Create)
		comGrp.POST("/update", admincom.Update)
		comGrp.POST("/remove", admincom.Remove)
		comGrp.POST("/onSale", admincom.OnSale)
		comGrp.POST("/offSale", admincom.OffSale)
		comGrp.GET("/content/:id", admincom.Content)
	}

	comspecGrp := adGrp.Group("/comspec")
	comspecGrp.Use(mdw.LoginAPIRequired())
	{
		comspecGrp.GET("/list", admincomspec.List)
		comspecGrp.POST("/create", admincomspec.Create)
		comspecGrp.POST("/valueCreate", admincomspec.ValueCreate)

	}

	freightGrp := adGrp.Group("/freight")
	freightGrp.Use(mdw.LoginAPIRequired())
	{
		freightGrp.POST("/list", adminfreight.List)
	}

	// 分类信息
	editorGrp := adGrp.Group("/editor")
	editorGrp.Use(mdw.LoginAPIRequired())
	{
		// 保存feed内容
		editorGrp.POST("/release", admineditor.Release)
		editorGrp.POST("/contentSave", admineditor.ContentSave)
		editorGrp.POST("/upload", admineditor.Upload)
	}

	// 分类信息
	comcateGrp := adGrp.Group("/comcate")
	comcateGrp.Use(mdw.LoginAPIRequired())
	{
		comcateGrp.GET("/list", admincomcate.List)
		comcateGrp.GET("/info", admincomcate.Info)
		comcateGrp.POST("/create", admincomcate.Create)
		comcateGrp.POST("/update", admincomcate.Update)
		comcateGrp.POST("/remove", admincomcate.Remove)
	}

	// 图片信息
	imageGrp := adGrp.Group("/image")
	imageGrp.Use(mdw.LoginAPIRequired())
	{
		imageGrp.GET("/list", image.List)
		imageGrp.POST("/add", image.Add)
	}
}

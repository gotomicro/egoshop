package router

import (
	"github.com/gin-gonic/gin"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/pkg/mus"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/admin/admincom"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/admin/admincomcate"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/admin/admincomment"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/admin/admincomspec"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/admin/admineditor"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/admin/adminfreight"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/admin/adminuser"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/admin/adminusergood"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/admin/auth"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/admin/image"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router/mdw"
)

func InitRouter() *gin.Engine {
	r := mus.Gin
	adminGrp(r) // 后台admin路由组
	viewGrp(r)  // 后台admin路由组
	return r
}

func viewGrp(r *gin.Engine) {
	r.Static("/static", "static")
	r.Static("/uploads", "uploads")
}

func adminGrp(r *gin.Engine) {
	//r.Use(mdw.JSONErr())

	// TODO 路由规范
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

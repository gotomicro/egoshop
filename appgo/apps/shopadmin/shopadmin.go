package main

import (
	"github.com/goecology/egoshop/appgo/apps/shopadmin/pkg/conf"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/pkg/mus"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/router"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/service"
	"github.com/goecology/muses"
	"github.com/goecology/muses/pkg/cache/redis"
	"github.com/goecology/muses/pkg/cmd"
	"github.com/goecology/muses/pkg/common"
	"github.com/goecology/muses/pkg/database/mysql"
	musgin "github.com/goecology/muses/pkg/server/gin"
	"github.com/goecology/muses/pkg/server/stat"
	"github.com/goecology/muses/pkg/session/ginsession"
)

func main() {
	app := muses.Container(
		cmd.Register,
		stat.Register,
		ginsession.Register,
		redis.Register,
		mysql.Register,
		musgin.Register,
	)
	app.SetRouter(router.InitRouter)
	app.PreRun(mus.Init, service.Init)
	err := conf.InitConfig(common.CmdConfigPath)
	if err != nil {
		panic(err)
	}
	err = app.Run()
	if err != nil {
		panic(err)
	}
}

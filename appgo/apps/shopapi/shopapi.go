package main

import (
	"github.com/goecology/egoshop/appgo/apps/shopapi/pkg/conf"
	"github.com/goecology/egoshop/appgo/apps/shopapi/pkg/mus"
	"github.com/goecology/egoshop/appgo/apps/shopapi/router"
	"github.com/goecology/egoshop/appgo/apps/shopapi/service"
	"github.com/goecology/muses"
	"github.com/goecology/muses/pkg/cache/redis"
	"github.com/goecology/muses/pkg/cmd"
	"github.com/goecology/muses/pkg/database/mysql"
	musgin "github.com/goecology/muses/pkg/server/gin"
	"github.com/goecology/muses/pkg/server/stat"
)

func main() {
	app := muses.Container(
		cmd.Register,
		stat.Register,
		redis.Register,
		mysql.Register,
		musgin.Register,
	)
	app.SetRouter(router.InitRouter)
	app.PreRun(mus.Init, conf.Init, service.Init)
	err := app.Run()
	if err != nil {
		panic(err)
	}
}

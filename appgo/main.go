package main

import (
	"github.com/goecology/egoshop/appgo/command"
	"github.com/goecology/egoshop/appgo/pkg/conf"
	"github.com/goecology/egoshop/appgo/pkg/mus"
	"github.com/goecology/egoshop/appgo/router"
	"github.com/goecology/egoshop/appgo/service"
	"github.com/goecology/muses"
	"github.com/goecology/muses/pkg/cache/redis"
	"github.com/goecology/muses/pkg/cmd"
	"github.com/goecology/muses/pkg/database/mysql"
	"github.com/goecology/muses/pkg/oss"
	musgin "github.com/goecology/muses/pkg/server/gin"
	"github.com/goecology/muses/pkg/server/stat"
	"github.com/goecology/muses/pkg/session/ginsession"
	"github.com/spf13/cobra"
)

func main() {
	app := muses.Container(
		cmd.Register,
		stat.Register,
		redis.Register,
		mysql.Register,
		musgin.Register,
		oss.Register,
		ginsession.Register,
	)
	app.SetRootCommand(func(cobraCommand *cobra.Command) {
		cobraCommand.AddCommand(command.InstallCmd)
	})
	app.SetStartCommand(func(cobraCommand *cobra.Command) {
		cobraCommand.PersistentFlags().StringVar(&command.Mode, "mode", "all", "设置启动模式")
	})
	app.SetGinRouter(router.InitRouter)
	app.PreRun(conf.Init, mus.Init, service.Init)
	err := app.Run()
	if err != nil {
		panic(err)
	}
}

package cmd

import (
	"fmt"
	"github.com/goecology/egoshop/appgo/apps/shopapi/service"
	"syscall"

	"github.com/fvbock/endless"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/goecology/egoshop/appgo/apps/shopapi/pkg/bootstrap"
	"github.com/goecology/egoshop/appgo/apps/shopapi/pkg/mus"
	"github.com/goecology/egoshop/appgo/apps/shopapi/router"
	"github.com/goecology/muses/pkg/server/gin"
)

var startCmd = &cobra.Command{
	Use:  "start",
	Long: `Starts egoshop server`,
	Run:  startFn,
}

func init() {
	startCmd.PersistentFlags().StringVarP(&bootstrap.Arg.CfgFile, "config", "c", "conf/conf.toml", "config file (default is $HOME/.cobra-example.yaml)")
	startCmd.PersistentFlags().BoolVarP(&bootstrap.Arg.Local, "local", "l", false, "local mode")
	RootCmd.AddCommand(startCmd)
	cobra.OnInitialize(initConfig)
}

func startFn(cmd *cobra.Command, args []string) {
	bootstrap.Init()
	service.Init()

	// 主服务器
	endless.DefaultReadTimeOut = gin.Config().Muses.Server.Gin.ReadTimeout.Duration
	endless.DefaultWriteTimeOut = gin.Config().Muses.Server.Gin.WriteTimeout.Duration
	endless.DefaultMaxHeaderBytes = 100000000000000
	server := endless.NewServer(gin.Config().Muses.Server.Gin.Addr, router.InitRouter())
	server.BeforeBegin = func(add string) {
		mus.Logger.Info(fmt.Sprintf("Actual pid is %d", syscall.Getpid()))
	}

	if err := server.ListenAndServe(); err != nil {
		mus.Logger.Error("Server err", zap.String("err", err.Error()))
	}
}

func initConfig() {
	viper.SetConfigFile(bootstrap.Arg.CfgFile)
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

package bootstrap

import (
	"github.com/goecology/egoshop/appgo/apps/shopadmin/pkg/conf"
	"github.com/goecology/egoshop/appgo/apps/shopadmin/pkg/mus"
)

var Arg arg

type arg struct {
	CfgFile string
	Local   bool
}

func Init() {
	err := conf.InitConfig(Arg.CfgFile)
	if err != nil {
		panic("config is error,err is" + err.Error())
	}
	mus.Init(Arg.CfgFile)
}

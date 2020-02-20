package command

import (
	"github.com/goecology/egoshop/appgo/command/install"
	"github.com/goecology/muses"
	"github.com/goecology/muses/pkg/cache/mixcache"
	mmysql "github.com/goecology/muses/pkg/database/mysql"
	"github.com/goecology/muses/pkg/oss"
	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:  "install",
	Long: `Show install information`,
	Run:  installCmd,
}

var ConfigPath string
var InstallMode string
var ClearMode bool

func init() {
	InstallCmd.PersistentFlags().StringVarP(&ConfigPath, "conf", "c", "conf/conf.toml", "conf path")
	InstallCmd.PersistentFlags().StringVarP(&InstallMode, "mode", "m", "all", "mode type")
	InstallCmd.PersistentFlags().BoolVar(&ClearMode, "clear", false, "dangerous mode")
}

func installCmd(cmd *cobra.Command, args []string) {
	app := muses.Container(
		mmysql.Register,
		oss.Register,
		mixcache.Register,
	)
	app.SetCfg(ConfigPath)
	app.PreRun(func() error {
		var err error
		if InstallMode == "all" || InstallMode == "create" {
			err = install.Create(ClearMode)
			if err != nil {
				return err
			}
		}
		if InstallMode == "all" || InstallMode == "mock" {
			err = install.Mock()
			if err != nil {
				return err
			}
		}
		return nil
	})
	err := app.Run()
	if err != nil {
		panic(err)
	}
}

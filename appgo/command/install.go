package command

import (
	"github.com/i2eco/egoshop/appgo/command/install"
	"github.com/i2eco/muses"
	"github.com/i2eco/muses/pkg/cache/mixcache"
	mmysql "github.com/i2eco/muses/pkg/database/mysql"
	"github.com/i2eco/muses/pkg/oss"
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
	app.SetPostRun(func() error {
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

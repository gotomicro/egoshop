package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/pflag"
	"os"
	"os/exec"
)

var (
	generator    string
	dbName       string
	mysqlHandler string
	modelPath    string
	daoPath      string
	outdaoPath   string
	appPath      string
	modulePath   string
)

func init() {
	pflag.StringVarP(&generator, "generator", "g", "", `指定generator仓库所在的路径`)
	pflag.StringVarP(&dbName, "db", "d", "egoshop", `指定数据库名`)
	pflag.StringVarP(&mysqlHandler, "mysql", "m", "", `指定存储(MySQL等)地址`)
	pflag.StringVarP(&modelPath, "model", "", "", `指定存储(MySQL等)地址`)
	pflag.StringVarP(&outdaoPath, "outdao", "", "", `指定存储(MySQL等)地址`)
	pflag.StringVarP(&daoPath, "dao", "", "", `指定存储(MySQL等)地址`)
	pflag.StringVarP(&appPath, "app", "", "", `指定存储(MySQL等)地址`)
	pflag.StringVarP(&modulePath, "module", "", "", `指定存储(MySQL等)地址`)
	pflag.Parse()
}

func main() {
	codeGen()
}

func codeGen() {
	fmt.Println(`generator, mysql, db, root, module------>`, generator, mysqlHandler, dbName, appPath, appPath)
	cmdStr := fmt.Sprintf(`
	cd %s && \
  	bin/generator new \
		--mysql '%s/information_schema' \
		--db '%s' \
		--out '%s' \
		--module='%s' \
		--model='%s' \
		--dao='%s' \
		--outdao='%s' \
		--debug='%s' \
		&& goimports -w %s && goimports -w %s`, generator, mysqlHandler, dbName, appPath, modulePath, modelPath, daoPath, outdaoPath, "true", outdaoPath, appPath)

	fmt.Println(cmdStr)
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	fmt.Println("[migrate] codegen done")
}

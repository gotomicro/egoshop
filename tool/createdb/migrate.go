package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/pflag"
)

var (
	dbName       string
	mysqlHandler string
)

func init() {
	pflag.StringVarP(&dbName, "db", "d", "egoshop", `指定数据库名`)
	pflag.StringVarP(&mysqlHandler, "mysql", "m", "", `指定存储(MySQL等)地址`)
	pflag.Parse()
}

func main() {
	db, e := gorm.Open(
		"mysql",
		fmt.Sprintf("%s/%s?charset=utf8&parseTime=True&loc=Local", mysqlHandler, dbName),
	)
	db = db.Debug()
	defer db.Close()

	if e != nil {
		fmt.Println("[migrate] conn fail", e)
		return
	}
	db.SingularTable(true)
	if db.Error != nil {
		fmt.Println("db.err", db.Error)
	}
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(Models...)
	if db.Error != nil {
		fmt.Println("db.err", db.Error)
	}
	fmt.Println("[migrate] migrate done")
}

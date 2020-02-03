package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goecology/muses"
	"github.com/goecology/muses/pkg/cache/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/pflag"
)

var (
	dbName       string
	mysqlHandler string
	endpoints    string
	key          string
	secret       string
	bucket       string
)

func init() {
	pflag.StringVarP(&dbName, "db", "d", "egoshop", `指定数据库名`)
	pflag.StringVarP(&mysqlHandler, "mysql", "m", "", `指定存储(MySQL等)地址`)
	pflag.StringVar(&endpoints, "endpoints", "m", `指定存储(MySQL等)地址`)
	pflag.StringVar(&key, "key", "", `指定存储(MySQL等)地址`)
	pflag.StringVar(&secret, "secret", "", `指定存储(MySQL等)地址`)
	pflag.StringVar(&bucket, "bucket", "", `指定存储(MySQL等)地址`)
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
	db.DropTable(Models...)
	if db.Error != nil {
		fmt.Println("db.err", db.Error)
		return
	}
	db.SingularTable(true)

	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(Models...)
	if db.Error != nil {
		fmt.Println("db.err", db.Error)
		return
	}

	if err := muses.Container(
		[]byte(`[muses.redis.egoshop]
        debug = true
        addr = "127.0.0.1:26479"
        network = "tcp"
        db = 0
        password = ""
        connectTimeout = "1s"
        readTimeout = "1s"
        writeTimeout = "1s"
        maxIdle = 5
        maxActive = 20
        idleTimeout = "60s"
        wait = false`),
		redis.Register,
	); err != nil {
		panic(err)
	}
	// 添加mock数据
	mock(db, endpoints, key, secret, bucket, redis.Caller("egoshop"))
	fmt.Println("[migrate] migrate done")
}

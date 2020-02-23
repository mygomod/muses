package main

import (
	"github.com/goecology/muses"
	"github.com/goecology/muses/pkg/app"
	"github.com/goecology/muses/pkg/cache/redis"
	"github.com/goecology/muses/pkg/database/mongo"
	"github.com/goecology/muses/pkg/database/mysql"
	"github.com/goecology/muses/pkg/logger"
	"github.com/goecology/muses/pkg/server/gin"
	"github.com/goecology/muses/pkg/server/stat"
	"github.com/goecology/muses/pkg/session/ginsession"
	"github.com/jinzhu/gorm"
)

var cfg = `
[muses.app]
	name = "muses"
	env = "dev"
	version = "1.0"
[muses.logger.system]
    debug = true
    level = "debug"
    path = "./system.log"
[muses.mysql.default]
    debug = true
    level = "panic"
    network = "tcp"
    dialect = "mysql"
    addr = "127.0.0.1:3306"
    username = "root"
    password = "root"
    db = "shop"
    charset = "utf8"
    parseTime = "True"
    loc = "Local"
    timeout = "1s"
    readTimeout = "1s"
    writeTimeout = "1s"
    maxOpenConns = 30
    maxIdleConns = 10
    connMaxLifetime = "300s"
[muses.mongo.default]
	debug = true
	url = "mongodb://127.0.0.1:27017/admin"
[muses.redis.default]
    debug = true
    addr = "127.0.0.1:6379"
    network = "tcp"
    db = 0
    password = ""
    connectTimeout = "1s"
    readTimeout = "1s"
    writeTimeout = "1s"
    maxIdle = 5
    maxActive = 20
    idleTimeout = "60s"
    wait = false
[muses.server.stat]
	addr = ":8100"
	writeTimeout = "1s"
	readTimeout = "1s"
[muses.server.gin]
    graceful = true
	mod = "debug"
    addr = ":10004"
    writeTimeout = "10s"
    readTimeout = "10s"
    maxHeaderBytes = 100000000000000
	enabledRecovery = true
	enabledLogger = true
	enabledMetric = true
[muses.session.gin]
    name = "mysession"
    size = 10
    network = "tcp"
    addr = "127.0.0.1:6379"
    pwd = ""
    keypairs = "secret"
`
var (
	Db *gorm.DB
)

func main() {
	if err := muses.Container(
		[]byte(cfg),
		app.Register,
		mysql.Register,
		mongo.Register,
		redis.Register,
		ginsession.Register,
		logger.Register,
		stat.Register,
		gin.Register,
	); err != nil {
		panic(err)
	}

	initCaller()
	type User struct {
		Uid  int
		Name string
	}
	u := User{}
	Db.Table("user").Where("uid=?", 1).Find(&u)
	select {}
}

func initCaller() {
	Db = mysql.Caller("default")
}

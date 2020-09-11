package main

import (
	"github.com/mygomod/muses"
	"github.com/mygomod/muses/pkg/database/mysql"
)

var cfg = `
[muses.mysql.ecology]
    debug = true
    level = "panic"
    network = "tcp"
    dialect = "mysql"
    addr = "127.0.0.1:3306"
    username = "root"
    password = "root"
    db = "ecology"
    charset = "utf8"
    parseTime = "True"
    loc = "Local"
    timeout = "1s"
    readTimeout = "1s"
    writeTimeout = "1s"
    maxOpenConns = 30
    maxIdleConns = 10
    connMaxLifetime = "300s"

`

type User struct {
	MemberId int
}

func main() {
	app := muses.Container(
		mysql.Register,
	)
	app.SetCfg([]byte(cfg))
	err := app.Run()
	if err != nil {
		panic(err)
	}
	u := User{}
	if mysql.Caller("ecology") != nil {
		mysql.Caller("ecology").Table("member").Where("member_id=?", 1).Find(&u)
	}
}

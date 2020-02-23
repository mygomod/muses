package main

import (
	"github.com/i2eco/muses"
	"github.com/i2eco/muses/pkg/cmd"
	"github.com/i2eco/muses/pkg/database/mysql"
)

type User struct {
	MemberId int
}

func main() {
	app := muses.Container(
		cmd.Register,
		mysql.Register,
	)
	err := app.Run()
	if err != nil {
		panic(err)
	}
	u := User{}
	if mysql.Caller("ecology") != nil {
		mysql.Caller("ecology").Table("member").Where("member_id=?", 1).Find(&u)
	}
}

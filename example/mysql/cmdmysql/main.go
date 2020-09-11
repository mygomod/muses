package main

import (
	"github.com/mygomod/muses"
	"github.com/mygomod/muses/pkg/cmd"
	"github.com/mygomod/muses/pkg/database/mysql"
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

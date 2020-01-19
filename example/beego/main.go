package main

import (
	"fmt"
	"github.com/goecology/muses"
	"github.com/goecology/muses/pkg/tpl/tplbeego"
)

var cfg = `
[muses.tmpl.beego]
    debug = true
`

func main() {
	if err := muses.Container(
		[]byte(cfg),
		tplbeego.Register,
	); err != nil {
		panic(err)
	}
	obj, err := tplbeego.Caller()
	if err != nil {
		fmt.Println("err", err)
		return
	}
	obj.Data["hello"] = "hello goecology"
	output, err := obj.RenderBytes()
	if err != nil {
		fmt.Println("err", err)
		return
	}

	fmt.Println(string(output))
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/i2eco/muses"
	mgin "github.com/i2eco/muses/pkg/server/gin"
	"net/http"
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
[muses.server.gin]
    graceful = true
	mod = "debug"
    addr = ":10004"
    writeTimeout = "10s"
    readTimeout = "10s"
    maxHeaderBytes = 100000000000000
	enabledRecovery = true
	enabledLogger = true

`

func main() {
	app := muses.Container(
		mgin.Register,
	)

	app.SetCfg([]byte(cfg))
	err := app.Run()
	if err != nil {
		panic(err)
	}

	obj := mgin.Caller()
	obj.GET("/hello", func(context *gin.Context) {
		fmt.Println("context.Handler()===>", context.Handler())
		context.String(200, "%s", "hello")
	})

	http.ListenAndServe(":8005", obj)
}

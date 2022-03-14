package main

import (
	"github.com/ergoapi/exgin"
	"github.com/ergoapi/zlog"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	cfg := zlog.Config{
		Simple: true,
	}
	zlog.InitZlog(&cfg)
}

func main() {
	g := exgin.Init(true)
	g.Use(exgin.ExLog())
	g.Use(exgin.ExRecovery())
	g.Use(exgin.ExCors())
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))
	g.Run()
}

package main

import (
	"os"

	"github.com/ergoapi/exgin"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
}

func main() {
	cfg := exgin.Config{
		Debug:   true,
		Pprof:   true,
		Metrics: true,
		Gops:    true,
	}
	g := exgin.Init(&cfg)
	g.Use(exgin.ExTraceID())
	g.Use(exgin.ExLog("/metrics"))
	g.Use(exgin.ExRecovery())
	g.GET("/", func(ctx *gin.Context) {
		exgin.GinsData(ctx, nil, nil)
	})
	g.Run(":9999")
}

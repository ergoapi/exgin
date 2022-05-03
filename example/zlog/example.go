package main

import (
	"github.com/ergoapi/exgin"
	"github.com/ergoapi/zlog"
)

func init() {
	cfg := zlog.Config{
		Simple: true,
	}
	zlog.InitZlog(&cfg)
}

func main() {
	cfg := exgin.Config{
		Debug:   true,
		Pprof:   true,
		Metrics: true,
	}
	g := exgin.Init(&cfg)
	g.Use(exgin.ExZLog())
	g.Use(exgin.ExZRecovery())
	g.Use(exgin.ExCors())
	g.Run()
}

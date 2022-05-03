package main

import (
	"os"

	"github.com/ergoapi/exgin"
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
	g.Use(exgin.ExLLog("/metrics"))
	g.Use(exgin.ExLRecovery())
	g.Run()
}

//  Copyright (c) 2021. The EFF Team Authors.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  See the License for the specific language governing permissions and
//  limitations under the License.

package exgin

import (
	"fmt"
	"time"

	"github.com/ergoapi/util/exnet"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/google/gops/agent"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	promNamespace = "exgin"
	promGinLabels = []string{
		"status_code",
		"path",
		"method",
	}
	promGinReqCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: promNamespace,
			Name:      "req_count",
			Help:      "gin server request count",
		}, promGinLabels,
	)
	promGinReqLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: promNamespace,
			Name:      "req_latency",
			Help:      "gin server request latency in seconds",
		}, promGinLabels,
	)
	// 默认慢请求时间 3s
	defaultGinSlowThreshold = time.Second * 3
)

type Config struct {
	Debug       bool
	Gops        bool
	GopsPath    string
	Pprof       bool
	PprofPath   string
	Cors        bool
	Metrics     bool
	MetricsPath string
}

func (c *Config) GinSet(r *gin.Engine) {
	gin.DisableConsoleColor()
	if c.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if c.Cors {
		r.Use(ExCors())
	}
	if c.Gops {
		if c.GopsPath == "" {
			c.GopsPath = "0.0.0.0:32388"
		}
		go agent.Listen(agent.Options{
			Addr:            c.GopsPath,
			ShutdownCleanup: true})
	}
	if c.Pprof {
		if c.PprofPath == "" {
			c.PprofPath = fmt.Sprintf("/hostdebug/%v/entry", exnet.LocalIPs()[0])
		}
		pprof.Register(r, c.PprofPath)
	}
	if c.Metrics {
		if c.MetricsPath == "" {
			c.MetricsPath = "/metrics"
		}
		r.GET(c.MetricsPath, gin.WrapH(promhttp.Handler()))
	}
}

// Init init gin engine
func Init(c *Config) *gin.Engine {
	r := gin.New()
	c.GinSet(r)
	return r
}

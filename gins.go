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
	"time"

	"github.com/ergoapi/util/zos"
	"github.com/gin-gonic/gin"
	"github.com/google/gops/agent"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

// Init init gin engine
func Init(debug bool) *gin.Engine {
	gin.DisableConsoleColor()
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if !zos.IsMacOS() {
		go agent.Listen(agent.Options{
			Addr:            "0.0.0.0:18848",
			ShutdownCleanup: true})
	}
	return gin.New()
}

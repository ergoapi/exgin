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
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/ergoapi/errors"
	ltrace "github.com/ergoapi/llog/hooks/trace"
	"github.com/ergoapi/util/exid"
	"github.com/ergoapi/util/ztime"
	"github.com/ergoapi/zlog/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// ExCors excors middleware
func ExCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, UPDATE, HEAD, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Access-Control-Request-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Max-Age", "3600")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// ExZLog exzlog middleware
func ExZLog(skip ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		host := Host(c)
		path := c.Request.URL.Path
		method := c.Request.Method
		ua := c.Request.UserAgent()
		query := c.Request.URL.RawQuery
		log := zlog.GetLogger().GetCtx(c.Request.Context())
		c.Next()
		for _, s := range skip {
			if s == path {
				return
			}
		}
		end := time.Now()
		latency := end.Sub(start)
		if len(query) == 0 {
			query = " - "
		}
		if latency > defaultGinSlowThreshold {
			log.Warn("api query slow", zap.String("path", path), zap.String("latency", time.Duration(latency).String()))
		}
		statuscode := c.Writer.Status()
		bodysize := c.Writer.Size()
		log.Info("api query", zap.String("path", path), zap.String("method", method), zap.String("host", host), zap.String("ua", ua), zap.String("query", query), zap.Int("status", statuscode), zap.Int("bodysize", bodysize), zap.String("latency", time.Duration(latency).String()))
		if len(c.Errors) > 0 || c.Writer.Status() >= 500 {
			log.Warn("api error query",
				zap.String("path", path),
				zap.String("method", method),
				zap.String("host", host),
				zap.String("ua", ua),
				zap.String("ip", RealIP(c)),
				zap.String("query", query),
				zap.Int("status", statuscode),
				zap.Int("bodysize", bodysize),
				zap.String("latency", time.Duration(latency).String()),
				zap.String("err", c.Errors.String()))
		} else {
			log.Info("api custom query",
				zap.String("path", path),
				zap.String("method", method),
				zap.String("host", host),
				zap.String("ua", ua),
				zap.String("ip", RealIP(c)),
				zap.String("query", query),
				zap.Int("status", statuscode),
				zap.Int("bodysize", bodysize),
				zap.String("latency", time.Duration(latency).String()))
		}

		// update prom
		labels := []string{fmt.Sprint(statuscode), path, method}
		promGinReqCount.WithLabelValues(labels...).Inc()
		promGinReqLatency.WithLabelValues(labels...).Observe(latency.Seconds())
	}
}

func ExZTraceID() gin.HandlerFunc {
	return func(g *gin.Context) {
		traceId := g.GetHeader("X-Trace-Id")
		if traceId == "" {
			traceId = exid.GenUUID()
			g.Header("X-Trace-Id", traceId)
		}
		ctx, _ := zlog.GetLogger().AddCtx(g.Request.Context(), zap.Any("traceId", traceId))
		g.Request = g.Request.WithContext(ctx)
		g.Next()
	}
}

// ExZRecovery zlog recovery
func ExZRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if res, ok := err.(errors.ErgoError); ok {
					GinsData(c, nil, fmt.Errorf(res.Message))
					c.Abort()
					return
				}
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zlog.GetLogger().Error("Recovery from brokenPipe", zap.String("path", c.Request.URL.Path), zap.Any("error", err), zap.String("request", string(httpRequest)))
					c.AbortWithStatusJSON(200, gin.H{
						"data":      nil,
						"message":   "请求broken",
						"timestamp": ztime.NowUnix(),
						"code":      10500,
					})
				} else {
					zlog.GetLogger().Error("Recovery from panic", zap.String("path", c.Request.URL.Path), zap.Any("error", err), zap.String("stack", string(debug.Stack())), zap.String("request", string(httpRequest)))
					c.AbortWithStatusJSON(200, gin.H{
						"data":      nil,
						"message":   "请求panic",
						"timestamp": ztime.NowUnix(),
						"code":      10500,
					})
				}
				return
			}
		}()
		c.Next()
	}
}

func ExLTraceID() gin.HandlerFunc {
	return func(g *gin.Context) {
		traceId := g.GetHeader("X-Trace-Id")
		if traceId == "" {
			traceId = exid.GenUUID()
			g.Header("X-Trace-Id", traceId)
		}
		logrus.AddHook(ltrace.NewTraceIdHook(traceId))
		g.Next()
	}
}

// ExLLog ex logrus middleware
func ExLLog(skip ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		host := Host(c)
		path := c.Request.URL.Path
		method := c.Request.Method
		ua := c.Request.UserAgent()
		query := c.Request.URL.RawQuery
		c.Next()
		for _, s := range skip {
			if strings.HasPrefix(path, s) {
				return
			}
		}
		end := time.Now()
		latency := end.Sub(start)
		if len(query) == 0 {
			query = " - "
		}
		if latency > defaultGinSlowThreshold {
			logrus.Warnf("[msg] api %v query %v", path, latency)
		}
		statuscode := c.Writer.Status()
		bodysize := c.Writer.Size()
		if len(c.Errors) > 0 || c.Writer.Status() >= 500 {
			logrus.Warnf("requestid %v =>  %v | %v | %v | %v | %v | %v | %v | %v  <= err: %v", statuscode, bodysize, RealIP(c), method, host, path, query, latency, ua, c.Errors.String())
		} else {
			logrus.Infof("requestid %v =>  %v | %v | %v | %v | %v | %v | %v | %v", statuscode, bodysize, RealIP(c), method, host, path, query, latency, ua)
		}
		// update prom
		labels := []string{fmt.Sprint(statuscode), path, method}
		promGinReqCount.WithLabelValues(labels...).Inc()
		promGinReqLatency.WithLabelValues(labels...).Observe(latency.Seconds())
	}
}

// ExLRecovery logrus recovery
func ExLRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if res, ok := err.(errors.ErgoError); ok {
					GinsData(c, nil, fmt.Errorf(res.Message))
					c.Abort()
					return
				}
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logrus.Errorf("Recovery from brokenPipe ---> path: %v, err: %v, request: %v", c.Request.URL.Path, err, string(httpRequest))
					c.AbortWithStatusJSON(200, gin.H{
						"data":      nil,
						"message":   "请求broken",
						"timestamp": ztime.NowUnix(),
						"code":      10500,
					})
				} else {
					logrus.Errorf("Recovery from panic ---> err: %v, request: %v, stack: %v", err, string(httpRequest), string(debug.Stack()))
					c.AbortWithStatusJSON(200, gin.H{
						"data":      nil,
						"message":   "请求panic",
						"timestamp": ztime.NowUnix(),
						"code":      10500,
					})
				}
				return
			}
		}()
		c.Next()
	}
}

func RealIP(c *gin.Context) string {
	xff := c.Writer.Header().Get("X-Forwarded-For")
	if xff == "" {
		return c.ClientIP()
	}
	return xff
}

func Host(c *gin.Context) string {
	h := c.Request.Host
	if h == "" {
		return c.Request.URL.Host
	}
	return h
}

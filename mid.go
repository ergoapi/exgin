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
	"net/http"

	"github.com/gin-gonic/gin"
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

// ExLog exlog middleware
// func ExLog() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		start := time.Now()
// 		path := c.Request.URL.Path
// 		query := c.Request.URL.RawQuery
// 		c.Next()
// 		end := time.Now()
// 		latency := end.Sub(start)
// 		if len(query) == 0 {
// 			query = " - "
// 		}
// 		if latency > time.Second*1 {
// 			zlog.Warn("[msg] api %v query %v", path, latency)
// 		}
// 		if len(c.Errors) > 0 || c.Writer.Status() >= 500 {
// 			msg := fmt.Sprintf("requestid %v => %v | %v | %v | %v | %v | %v <= err: %v", rid.GetRID(c), color.SRed("%v", c.Writer.Status()), c.ClientIP(), c.Request.Method, path, query, latency, c.Errors.String())
// 			zlog.Warn(msg)
// 			go file.Writefile(fmt.Sprintf("/tmp/%v.errreq.txt", ztime.NowDay()), msg)
// 		} else {
// 			zlog.Info("requestid %v => %v | %v | %v | %v | %v | %v ", rid.GetRID(c), color.SGreen("%v", c.Writer.Status()), c.ClientIP(), c.Request.Method, path, query, latency)
// 		}
// 	}
// }

// Exrecovery recovery
// func Exrecovery() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				if res, ok := err.(gerr.GinsError); ok {
// 					engine.GinsData(c, nil, fmt.Errorf(res.Message))
// 					c.Abort()
// 					return
// 				}
// 				var brokenPipe bool
// 				if ne, ok := err.(*net.OpError); ok {
// 					if se, ok := ne.Err.(*os.SyscallError); ok {
// 						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
// 							brokenPipe = true
// 						}
// 					}
// 				}

// 				httpRequest, _ := httputil.DumpRequest(c.Request, false)
// 				if brokenPipe {
// 					zlog.Error("Recovery from brokenPipe ---> path: %v, err: %v, request: %v", c.Request.URL.Path, err, string(httpRequest))
// 					c.AbortWithStatusJSON(200, gin.H{
// 						"data":      nil,
// 						"message":   "请求broken",
// 						"timestamp": ztime.NowUnix(),
// 						"code":      10500,
// 					})
// 				} else {
// 					zlog.Error("Recovery from panic ---> err: %v, request: %v, stack: %v", err, string(httpRequest), string(debug.Stack()))
// 					c.AbortWithStatusJSON(200, gin.H{
// 						"data":      nil,
// 						"message":   "请求panic",
// 						"timestamp": ztime.NowUnix(),
// 						"code":      10500,
// 					})
// 				}
// 				return
// 			}
// 		}()
// 		c.Next()
// 	}
// }
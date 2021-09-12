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

	"github.com/ergoapi/errors"
	"github.com/ergoapi/util/ztime"
	"github.com/gin-gonic/gin"
)

// done done
func respdone(data interface{}) gin.H {
	return gin.H{
		"data":      data,
		"message":   "请求成功",
		"timestamp": ztime.NowUnix(),
		"code":      200,
	}
}

// error error
func resperror(code int, data interface{}) gin.H {
	return gin.H{
		"data":      nil,
		"message":   data,
		"timestamp": ztime.NowUnix(),
		"code":      code,
	}
}

func renderMessage(c *gin.Context, v interface{}) {
	if v == nil {
		c.JSON(200, respdone(nil))
		return
	}

	switch t := v.(type) {
	case string:
		c.JSON(200, resperror(10400, t))
	case error:
		c.JSON(200, resperror(10400, t.Error()))
	}
}

func GinsData(c *gin.Context, data interface{}, err error) {
	if err == nil {
		c.JSON(200, respdone(data))
		return
	}

	renderMessage(c, err.Error())
}

func GinsAbort(c *gin.Context, msg string, args ...interface{}) {
	c.AbortWithStatusJSON(200, resperror(10400, fmt.Sprintf(msg, args...)))
}

func GinsAbortWithCode(c *gin.Context, respcode int, msg string, args ...interface{}) {
	c.AbortWithStatusJSON(200, resperror(respcode, fmt.Sprintf(msg, args...)))
}

func GinsCustomResp(c *gin.Context, obj interface{}) {
	c.JSON(200, obj)
}

func Bind(c *gin.Context, ptr interface{}) {
	err := c.ShouldBindJSON(ptr)
	if err != nil {
		errors.Bomb("参数不合法: %v", err)
	}
}

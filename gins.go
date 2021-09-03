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
	"github.com/gin-gonic/gin"
)

// Init init gin engine
func Init(debug bool) *gin.Engine {
	// status := false
	// if len(debug) == 1 {
	// 	status = debug[0]
	// }
	// if status || zos.IsMacOS() {
	// 	gin.SetMode(gin.DebugMode)
	// } else {
	// 	gin.SetMode(gin.ReleaseMode)
	// }
	// if zos.IsLinux() {
	// 	go agent.Listen(agent.Options{
	// 		Addr:            "0.0.0.0:8848",
	// 		ShutdownCleanup: true})
	// }
	gin.DisableConsoleColor()
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	return gin.New()
}

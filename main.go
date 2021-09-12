// Copyright 2021 KMU KCC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main runs the API server of the Buddy System.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/config"
	"github.com/kmu-kcc/buddy-backend/web/api/v1/activity"
	"github.com/kmu-kcc/buddy-backend/web/api/v1/fee"
	"github.com/kmu-kcc/buddy-backend/web/api/v1/member"
)

func main() {
	parser := argparse.NewParser("buddy", "API server of the Buddy System")

	// parse port number from command line arguments
	//
	// See https://github.com/akamensky/argparse#readme
	//
	//
	// NOTE:
	//
	// argparse is redundant due to the `flag` package in the standard library.
	// This would be removed in v1.1.0.
	port := parser.Int("p", "port", &argparse.Options{Required: true, Help: "Port to run the server"})

	if err := parser.Parse(os.Args); err != nil {
		log.Fatalln(parser.Usage(err))
	}

	gin.SetMode(gin.ReleaseMode)

	engine := gin.Default()

	engine.Use(cors.New(config.CORSConfig))

	api := engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			members := v1.Group("/member")
			{
				members.POST("/signin", member.SignIn())
				members.POST("/signup", member.SignUp())
				members.GET("/signups", member.SignUps())
				members.PUT("/approve", member.Approve())
				members.DELETE("/delete", member.Delete())
				members.PUT("/exit", member.Exit())
				members.GET("/exits", member.Exits())
				members.POST("/my", member.My())
				members.GET("/search", member.Search())
				members.PUT("/update", member.Update())
				members.GET("/active", member.Active())
				members.PUT("/activate", member.Activate())
				members.GET("/graduates", member.Graduates())
				members.PUT("/updaterole", member.UpdateRole())
			}
			activities := v1.Group("/activity")
			{
				activities.POST("/create", activity.Create())
				activities.GET("/search", activity.Search())
				activities.GET("/private", activity.Private())
				activities.PUT("/update", activity.Update())
				activities.DELETE("/delete", activity.Delete())
				activities.POST("/upload", activity.Upload())
				activities.POST("/download", activity.Download())
				activities.POST("/deletefile", activity.DeleteFile())
			}
			fees := v1.Group("/fee")
			{
				fees.POST("/create", fee.Create())
				fees.POST("/amount", fee.Amount())
				fees.POST("/payers", fee.Payers())
				fees.POST("/deptors", fee.Deptors())
				fees.POST("/search", fee.Search())
				fees.POST("/pay", fee.Pay())
				fees.POST("/deposit", fee.Deposit())
				fees.POST("/exempt", fee.Exempt())
			}
		}
	}

	log.Fatalln(engine.Run(fmt.Sprintf(":%d", *port)))
}

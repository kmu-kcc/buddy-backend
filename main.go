// Package main runs the API server of the Buddy System.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/web/api/v1/activity"
	"github.com/kmu-kcc/buddy-backend/web/api/v1/fee"
	"github.com/kmu-kcc/buddy-backend/web/api/v1/member"
)

func main() {
	parser := argparse.NewParser("buddy", "API server of the Buddy System")

	// parse port number from command line arguments
	//
	// See https://github.com/akamensky/argparse#readme
	port := parser.Int("p", "port", &argparse.Options{Required: true, Help: "Port to run the server"})

	if err := parser.Parse(os.Args); err != nil {
		log.Fatalln(parser.Usage(err))
	}

	gin.SetMode(gin.DebugMode)

	engine := gin.Default()

	api := engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			mgroup := v1.Group("/member")
			{
				mgroup.POST("/signup", member.SignUp())
			}
			agroup := v1.Group("/activity")
			{
				agroup.POST("/applyp", activity.ApplyP())
				agroup.GET("/papplies", activity.Papplies())
				agroup.POST("/approvep", activity.ApproveP())
				agroup.POST("/rejectp", activity.RejectP())
				agroup.POST("/cancelp", activity.CancelP())

			}
			fgroup := v1.Group("/fee")
			{
				fgroup.GET("/dones", fee.Dones())
				fgroup.GET("/yets", fee.Yets())
				fgroup.GET("/all", fee.All())
			}
		}
	}

	log.Fatalln(engine.Run(fmt.Sprintf(":%d", *port)))
}

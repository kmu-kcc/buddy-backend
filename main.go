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

	gin.SetMode(gin.ReleaseMode)

	engine := gin.Default()

	api := engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			mgroup := v1.Group("/member")
			{
				_ = mgroup
			}

			agroup := v1.Group("/activity")
			{
				agroup.POST("/applyc", activity.ApplyC())
				agroup.POST("/cancelc", activity.CancelC())
				agroup.GET("/capplies", activity.Capplies())
				agroup.POST("/approvec", activity.ApproveC())
				agroup.POST("/rejectc", activity.RejectC())
			}

			fgroup := v1.Group("/fee")
			{
				fgroup.POST("/approve", fee.Approve())
				fgroup.POST("/reject", fee.Reject())
				fgroup.POST("/deposit", fee.Deposit())
			}

		}
	}

	log.Fatalln(engine.Run(fmt.Sprintf(":%d", *port)))
}

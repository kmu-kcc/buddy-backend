// Package main runs the API server of the Buddy System.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/gin-gonic/gin"
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

	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	engine := gin.Default()

	log.Fatalln(engine.Run(fmt.Sprintf(":%d", *port)))
}

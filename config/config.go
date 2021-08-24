package config

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
)

var (
	MongoURI     = os.Getenv("MONGO_URI")
	AccessSecret = os.Getenv("ACCESS_SECRET")
	CORSConfig   = cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           6 * time.Hour,
	}
)

package config

import "os"

var (
	MongoURI     = os.Getenv("MONGO_URI")
	AccessSecret = os.Getenv("ACCESS_SECRET")
)

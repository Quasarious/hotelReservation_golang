package my_conf

import (
	"os"
)

var DBNAME string
var DBURL string

func LoadEnv() error {
	os.Setenv("MONGO_DB_NAME", "hotel-reservation")
	os.Setenv("MONGO_DB_URL", "mongodb://localhost:27017")
	os.Setenv("JWT_SECRET", "1882UtaCat")
	os.Setenv("HTTP_LISTEN_ADDRESS", ":3000")

	DBNAME = os.Getenv("MONGO_DB_NAME")
	DBURL = os.Getenv("MONGO_DB_URL")
	return nil
}

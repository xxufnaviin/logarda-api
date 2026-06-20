package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var ENCRYPTION_KEY string
var POSTGRES_HOST string
var POSTGRES_DATABASE string
var POSTGRES_USER string
var POSTGRES_PASSWORD string
var POSTGRES_URL string
var REDIS_HOST string
var REDIS_PORT string
var REDIS_URL string

func LoadSecrets() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	ENCRYPTION_KEY = os.Getenv("ENCRYPTION_KEY")
	POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	POSTGRES_DATABASE = os.Getenv("POSTGRES_DATABASE")
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	REDIS_HOST = os.Getenv("REDIS_HOST")
	REDIS_PORT = os.Getenv("REDIS_PORT")

	POSTGRES_URL = fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		POSTGRES_USER,
		POSTGRES_PASSWORD,
		POSTGRES_HOST,
		POSTGRES_DATABASE)
	
	REDIS_URL = fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT)
}

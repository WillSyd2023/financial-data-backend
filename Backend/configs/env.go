package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURL() string {
	// For testing, Load must be given "../.env"
	// For running, Load must be given nothing
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("MONGOURL")
}

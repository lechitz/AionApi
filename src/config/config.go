package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var (
	Port             = 5000
	StringConnection = ""
)

func LoadEnv() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	Port, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		Port = 9000
	}

	StringConnection = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	fmt.Println("Configurations loaded successfully")
	fmt.Printf(" - Connection: %s\n", StringConnection)
	fmt.Printf(" - Port: %d\n\n", Port)
}

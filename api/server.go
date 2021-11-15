package api

import (
	"fmt"
	"log"
	"os"

	"github.com/doka-guide/api/api/controllers"
	"github.com/doka-guide/api/api/seed"

	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Не могу получить доступ к файлу '.env': %v", err)
	} else {
		fmt.Println("Значения из файла '.env' получены.")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	seed.Load(server.DB)
	server.Run(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))
}

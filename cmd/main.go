package main

import (
	"fmt"
	"log"
	"os"
	"wechat-back/internals/routes"
	"wechat-back/internals/server"

	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load("db.env", ".env")
	if err != nil {
		fmt.Println("-------------- You will need to load .env files before starting server. --------------")
		fmt.Println("Please reffer to documentation at github.com/JENoyola/wechat-backend")
		log.Panic(err)
		return
	}

}

func main() {

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "2565"
	}

	config := server.ServeConfig{
		PORT:        PORT,
		HANDLER:     routes.ServerRoutes(),
		IDLE:        1,
		WRITE:       1,
		READHEADER:  1,
		READTIMEOUT: 1,
		TLSC:        os.Getenv("TLS_CERT"),
		TLSK:        os.Getenv("TSL_KEY"),
		ENV:         os.Getenv("ENV"),
		API_VERSION: os.Getenv("API_VERSION"),
	}

	log.Fatal(server.StartServer(config))

}

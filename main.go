package main

import (
	"api-golang-compro/config"
	"api-golang-compro/routes"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// for load godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}


	// database connection
	db := config.ConnectDatabase()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// router
	r := routes.SetupRouter(db)
	r.Run(":8001")
}
package config

import (
	"api-golang-compro/models"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	// koneksi dari ENV
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_NAME")

	// koneksi untuk docker
	// username := "api-golang-user"
	// password := "mysecretpassword"
	// host := "tcp(db:3306)"
	// database := "api-golang-compro"

	// koneksi untuk local
	// username := "root"
	// password := ""
	// host := "tcp(localhost:3306)"
	// database := "api-golang-compro"

	dsn := fmt.Sprintf("%v:%v@%v/%v?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&models.Career{}, &models.Category{}, &models.Position{}, &models.User{}, &models.Category_home{}, models.Form{}, models.Home{}, models.Navbar{}, models.Portfolio{}, models.Service{}, models.PortfolioHepytech{}, models.LogActivity{}, models.Logo{}, models.Impact{})

	return db
}

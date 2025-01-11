package models

import (
	"fmt"
	"log"

	"github.com/beego/beego/v2/server/web"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Fetch database configurations from app.conf
	driver, _ := web.AppConfig.String("db::driver")
	user, _ := web.AppConfig.String("db::user")
	password, _ := web.AppConfig.String("db::password")
	host, _ := web.AppConfig.String("db::host")
	port, _ := web.AppConfig.String("db::port")
	name, _ := web.AppConfig.String("db::name")

	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, name, port)

	// Initialize the database connection
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully!")
}

package conf

import (
    "fmt"
    "log"

    "github.com/beego/beego/v2/server/web"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
    host, err := web.AppConfig.String("db::host")
    if err != nil {
        log.Fatalf("Error reading db::host: %v", err)
    }
    user, err := web.AppConfig.String("db::user")
    if err != nil {
        log.Fatalf("Error reading db::user: %v", err)
    }
    password, err := web.AppConfig.String("db::password")
    if err != nil {
        log.Fatalf("Error reading db::password: %v", err)
    }
    dbname, err := web.AppConfig.String("db::dbname")
    if err != nil {
        log.Fatalf("Error reading db::dbname: %v", err)
    }
    port, err := web.AppConfig.String("db::port")
    if err != nil {
        log.Fatalf("Error reading db::port: %v", err)
    }

    fmt.Printf("host: %s, user: %s, password: %s, dbname: %s, port: %s\n", host, user, password, dbname, port)

    if host == "" || user == "" || password == "" || dbname == "" || port == "" {
        log.Fatalf("Database configuration values are missing")
    }

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        host, user, password, dbname, port)

    fmt.Println("DSN:", dsn) // Debug print DSN

    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }
    fmt.Println("Database connection initialized")
}
// // rental_api/conf/config.go
// package conf

// import (
//     "fmt"
//     "log"

//     "github.com/beego/beego/v2/server/web"
//     "gorm.io/driver/postgres"
//     "gorm.io/gorm"
// )

// var DB *gorm.DB

// // InitDB initializes the database connection
// func InitDB() {
//     host, _ := web.AppConfig.String("db::host")
//     user, _ := web.AppConfig.String("db::user")
//     password, _ := web.AppConfig.String("db::password")
//     dbname, _ := web.AppConfig.String("db::dbname")
//     port, _ := web.AppConfig.String("db::port")

//     dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
//         host, user, password, dbname, port)

//     var err error
//     DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
//     if err != nil {
//         log.Fatalf("Could not connect to the database: %v", err)
//     }
//     fmt.Println("Database connection initialized")
// }
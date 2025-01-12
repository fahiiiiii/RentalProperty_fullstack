//  File: conf/db.go
package conf

import (
    "fmt"
    "gopkg.in/gcfg.v1"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Config struct {
    Database struct {
        Host     string
        Port     string
        User     string
        Password string
        DBName   string
    }
}

var (
    DB     *gorm.DB
    AppConfig Config
)

// LoadConfig reads the configuration file
func LoadConfig() error {
    err := gcfg.ReadFileInto(&AppConfig, "conf/app.conf")
    if err != nil {
        return fmt.Errorf("failed to read config: %v", err)
    }
    return nil
}

// InitDB initializes the database connection
func InitDB() error {
    if err := LoadConfig(); err != nil {
        return err
    }

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        AppConfig.Database.Host,
        AppConfig.Database.Port,
        AppConfig.Database.User,
        AppConfig.Database.Password,
        AppConfig.Database.DBName,
    )

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    return err
}
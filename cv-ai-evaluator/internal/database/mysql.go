package database

import (
	"fmt"
	"log"
	"time"

	"cv-ai-evaluator/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
    var err error
    
    dsn := cfg.GetDSN()
    
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().Local()
        },
    })

    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    // Get underlying SQL DB untuk konfigurasi connection pool
    sqlDB, err := DB.DB()
    if err != nil {
        return fmt.Errorf("failed to get database instance: %w", err)
    }

    // Connection pool settings
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    log.Println("Database connected successfully!")
    return nil
}

func CloseDB() error {
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

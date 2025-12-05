package db

import (
	"fmt"
	"log"
	"os"

	"github.com/LimerDev/worklog/internal/config"
	"github.com/LimerDev/worklog/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	// Use config database settings, with environment variable overrides
	host := getEnv("DB_HOST", cfg.Database.Host)
	if host == "" {
		host = "localhost"
	}
	port := getEnv("DB_PORT", cfg.Database.Port)
	if port == "" {
		port = "5432"
	}
	user := getEnv("DB_USER", cfg.Database.User)
	if user == "" {
		user = "worklog"
	}
	password := getEnv("DB_PASSWORD", cfg.Database.Password)
	if password == "" {
		password = "worklog"
	}
	dbname := getEnv("DB_NAME", cfg.Database.Name)
	if dbname == "" {
		dbname = "worklog"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.Customer{},
		&models.Project{},
		&models.Consultant{},
		&models.TimeEntry{},
	)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

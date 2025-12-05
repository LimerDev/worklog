package database

import (
	"fmt"

	"github.com/LimerDev/worklog/internal/config"
	"github.com/LimerDev/worklog/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	// Read database configuration from config
	host := cfg.Database.Host
	port := cfg.Database.Port
	user := cfg.Database.User
	password := cfg.Database.Password
	dbname := cfg.Database.Name

	// Require all database configuration values
	if host == "" {
		return fmt.Errorf("database.host is required in config file (~/.worklog/config.json) or WORKLOG_DATABASE_HOST environment variable")
	}
	if port == "" {
		return fmt.Errorf("database.port is required in config file (~/.worklog/config.json) or WORKLOG_DATABASE_PORT environment variable")
	}
	if user == "" {
		return fmt.Errorf("database.user is required in config file (~/.worklog/config.json) or WORKLOG_DATABASE_USER environment variable")
	}
	if password == "" {
		return fmt.Errorf("database.password is required in config file (~/.worklog/config.json) or WORKLOG_DATABASE_PASSWORD environment variable")
	}
	if dbname == "" {
		return fmt.Errorf("database.name is required in config file (~/.worklog/config.json) or WORKLOG_DATABASE_NAME environment variable")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

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

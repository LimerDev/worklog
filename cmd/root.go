package cmd

import (
	"fmt"
	"os"

	"github.com/LimerDev/worklog/internal/config"
	db "github.com/LimerDev/worklog/internal/database"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "worklog",
	Short: "Worklog - A simple time reporting app for consultant hours",
	Long:  `Worklog is a CLI tool for registering and reporting consultant hours.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initDB)
}

func initConfig() {
	if err := config.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}
}

func initDB() {
	cfg, err := config.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read configuration: %v\n", err)
		os.Exit(1)
	}

	if err := db.Connect(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err := db.AutoMigrate(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to migrate database: %v\n", err)
		os.Exit(1)
	}
}

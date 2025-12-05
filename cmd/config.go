package cmd

import (
	"fmt"

	"github.com/LimerDev/worklog/internal/config"
	"github.com/spf13/cobra"
)

var (
	configConsultant string
	configClient     string
	configProject    string
	configRate       float64
	configDBHost     string
	configDBPort     string
	configDBUser     string
	configDBPassword string
	configDBName     string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage default values for consultant, client, and project",
	Long:  "Set or view default values to speed up time entry registration",
	RunE:  runConfigShow,
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set default values",
	Long:  "Set default consultant, client, project, and/or hourly rate",
	RunE:  runConfigSet,
}

var configClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all default values",
	Long:  "Remove all saved default values",
	RunE:  runConfigClear,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configClearCmd)

	configSetCmd.Flags().StringVarP(&configConsultant, "consultant", "n", "", "Default consultant name")
	configSetCmd.Flags().StringVarP(&configClient, "client", "c", "", "Default client name")
	configSetCmd.Flags().StringVarP(&configProject, "project", "p", "", "Default project name")
	configSetCmd.Flags().Float64VarP(&configRate, "rate", "r", 0, "Default hourly rate")
	configSetCmd.Flags().StringVar(&configDBHost, "db-host", "", "Database host")
	configSetCmd.Flags().StringVar(&configDBPort, "db-port", "", "Database port")
	configSetCmd.Flags().StringVar(&configDBUser, "db-user", "", "Database user")
	configSetCmd.Flags().StringVar(&configDBPassword, "db-password", "", "Database password")
	configSetCmd.Flags().StringVar(&configDBName, "db-name", "", "Database name")
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Printf("Configuration:\n\n")

	if cfg.DefaultConsultant == "" && cfg.DefaultClient == "" && cfg.DefaultProject == "" && cfg.DefaultRate == 0 {
		fmt.Println("No defaults configured yet.")
		fmt.Println("\nSet defaults with: worklog config set -n CONSULTANT -c CLIENT -p PROJECT -r RATE")
		return nil
	}

	if cfg.DefaultConsultant != "" {
		fmt.Printf("Default Consultant: %s\n", cfg.DefaultConsultant)
	}
	if cfg.DefaultClient != "" {
		fmt.Printf("Default Client: %s\n", cfg.DefaultClient)
	}
	if cfg.DefaultProject != "" {
		fmt.Printf("Default Project: %s\n", cfg.DefaultProject)
	}
	if cfg.DefaultRate > 0 {
		fmt.Printf("Default Hourly Rate: %.2f kr/h\n", cfg.DefaultRate)
	}

	fmt.Println("\nDatabase Configuration:")
	if cfg.Database.Host != "" {
		fmt.Printf("  Host: %s\n", cfg.Database.Host)
	}
	if cfg.Database.Port != "" {
		fmt.Printf("  Port: %s\n", cfg.Database.Port)
	}
	if cfg.Database.User != "" {
		fmt.Printf("  User: %s\n", cfg.Database.User)
	}
	if cfg.Database.Name != "" {
		fmt.Printf("  Database: %s\n", cfg.Database.Name)
	}

	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	if configConsultant == "" && configClient == "" && configProject == "" && configRate == 0 &&
		configDBHost == "" && configDBPort == "" && configDBUser == "" && configDBPassword == "" && configDBName == "" {
		return fmt.Errorf("you must specify at least one value")
	}

	if err := config.SaveDefaults(configConsultant, configClient, configProject, configRate); err != nil {
		return err
	}

	cfg, err := config.Get()
	if err != nil {
		return err
	}

	fmt.Println("  Configuration saved!")
	fmt.Printf("  Consultant: %s\n", cfg.DefaultConsultant)
	fmt.Printf("  Client: %s\n", cfg.DefaultClient)
	fmt.Printf("  Project: %s\n", cfg.DefaultProject)
	fmt.Printf("  Hourly Rate: %.2f kr/h\n", cfg.DefaultRate)

	return nil
}

func runConfigClear(cmd *cobra.Command, args []string) error {
	if err := config.ClearDefaults(); err != nil {
		return err
	}

	fmt.Println("✓ Configuration cleared!")

	return nil
}

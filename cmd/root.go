package cmd

import (
	"fmt"
	"os"

	"github.com/LimerDev/worklog/internal/config"
	db "github.com/LimerDev/worklog/internal/database"
	"github.com/LimerDev/worklog/internal/i18n"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "worklog",
	Short: "", // Set after i18n initialization
	Long:  "", // Set after i18n initialization
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Initialize config and i18n early so command descriptions can be localized
	initConfig()
	initI18n()

	// Now localize all commands after i18n is ready
	rootCmd.Short = i18n.T(i18n.KeyRootShort)
	rootCmd.Long = i18n.T(i18n.KeyRootLong)

	for _, c := range rootCmd.Commands() {
		localizeCommand(c)
	}

	rootCmd.PersistentPreRunE = persistentPreRun
	initialized = true
}

var initialized bool

func persistentPreRun(cmd *cobra.Command, args []string) error {
	// Initialize database for non-help commands
	if cmd.Name() != "help" && !cmd.Flags().Changed("help") {
		initDB()
	}

	return nil
}

func initConfig() {
	if err := config.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}
}

func initI18n() {
	cfg, err := config.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read configuration: %v\n", err)
		os.Exit(1)
	}

	if err := i18n.Initialize(cfg.Language); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize i18n: %v\n", err)
		os.Exit(1)
	}
}

func initDB() {
	cfg, err := config.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, i18n.T(i18n.KeyErrReadConfig)+": %v\n", err)
		os.Exit(1)
	}

	if err := db.Connect(cfg); err != nil {
		fmt.Fprintf(os.Stderr, i18n.T(i18n.KeyErrInitDatabase)+": %v\n", err)
		os.Exit(1)
	}

	if err := db.AutoMigrate(); err != nil {
		fmt.Fprintf(os.Stderr, i18n.T(i18n.KeyErrDatabaseMigrate)+": %v\n", err)
		os.Exit(1)
	}
}

func localizeCommand(cmd *cobra.Command) {
	switch cmd.Name() {
	case "add":
		localizeAddCommand()
	case "get":
		localizeGetCommand()
	case "export":
		localizeExportCommand()
	case "config":
		localizeConfigCommand()
	}
}

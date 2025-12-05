package cmd

import (
	"fmt"
	"os"

	"github.com/LimerDev/worklog/internal/config"
	"github.com/LimerDev/worklog/internal/db"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "worklog",
	Short: "Worklog - En enkel tidsrapporteringsapp för konsulttimmar",
	Long:  `Worklog är ett CLI-verktyg för att registrera och rapportera konsulttimmar.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initDB)
}

func initDB() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Kunde inte ladda konfiguration: %v\n", err)
		os.Exit(1)
	}

	if err := db.Connect(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Kunde inte ansluta till databasen: %v\n", err)
		os.Exit(1)
	}

	if err := db.AutoMigrate(); err != nil {
		fmt.Fprintf(os.Stderr, "Kunde inte migrera databasen: %v\n", err)
		os.Exit(1)
	}
}

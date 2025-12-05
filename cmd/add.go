package cmd

import (
	"fmt"
	"time"

	"github.com/LimerDev/worklog/internal/config"
	"github.com/LimerDev/worklog/internal/database"
	"github.com/LimerDev/worklog/internal/models"
	"github.com/spf13/cobra"
)

var (
	hours       float64
	description string
	project     string
	client      string
	consultant  string
	hourlyRate  float64
	date        string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Lägg till en ny tidsregistrering",
	Long:  `Registrera en ny tidsregistrering med timmar, beskrivning, projekt och kund.`,
	RunE:  runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().Float64VarP(&hours, "hours", "t", 0, "Antal timmar (obligatorisk)")
	addCmd.Flags().StringVarP(&description, "description", "d", "", "Beskrivning av arbetet (obligatorisk)")
	addCmd.Flags().StringVarP(&project, "project", "p", "", "Projektnamn (använder default om inte angiven)")
	addCmd.Flags().StringVarP(&client, "client", "c", "", "Kundnamn (använder default om inte angiven)")
	addCmd.Flags().StringVarP(&consultant, "consultant", "n", "", "Konsultnamn (använder default om inte angiven)")
	addCmd.Flags().Float64VarP(&hourlyRate, "rate", "r", 0, "Timpris (använder default om inte angivet)")
	addCmd.Flags().StringVarP(&date, "date", "D", "", "Datum (YYYY-MM-DD, standard: idag)")

	addCmd.MarkFlagRequired("hours")
	addCmd.MarkFlagRequired("description")
}

func runAdd(cmd *cobra.Command, args []string) error {
	// Get configuration
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("kunde inte läsa konfiguration: %w", err)
	}

	// Use defaults if not provided
	if consultant == "" {
		consultant = cfg.DefaultConsultant
	}
	if client == "" {
		client = cfg.DefaultClient
	}
	if project == "" {
		project = cfg.DefaultProject
	}
	if hourlyRate == 0 {
		hourlyRate = cfg.DefaultRate
	}

	// Validate required fields
	if consultant == "" {
		return fmt.Errorf("konsult krävs (-n KONSULT eller `worklog config set -n KONSULT`)")
	}
	if client == "" {
		return fmt.Errorf("kund krävs (-c KUND eller `worklog config set -c KUND`)")
	}
	if project == "" {
		return fmt.Errorf("projekt krävs (-p PROJEKT eller `worklog config set -p PROJEKT`)")
	}
	if hourlyRate <= 0 {
		return fmt.Errorf("timpris krävs (-r TIMPRIS eller `worklog config set -r TIMPRIS`)")
	}

	var entryDate time.Time

	if date == "" {
		entryDate = time.Now()
	} else {
		entryDate, err = time.Parse("2006-01-02", date)
		if err != nil {
			return fmt.Errorf("ogiltigt datumformat, använd YYYY-MM-DD: %w", err)
		}
	}

	if hours <= 0 {
		return fmt.Errorf("timmar måste vara större än 0")
	}

	repo := database.NewRepository()

	// Get or create consultant
	consultantObj, err := repo.GetOrCreateConsultant(consultant)
	if err != nil {
		return fmt.Errorf("kunde inte hämta/skapa konsult: %w", err)
	}

	// Get or create customer
	customerObj, err := repo.GetOrCreateCustomer(client)
	if err != nil {
		return fmt.Errorf("kunde inte hämta/skapa kund: %w", err)
	}

	// Get or create project for this customer
	projectObj, err := repo.GetOrCreateProject(project, customerObj.ID)
	if err != nil {
		return fmt.Errorf("kunde inte hämta/skapa projekt: %w", err)
	}

	// Create time entry
	entry := &models.TimeEntry{
		Date:         entryDate,
		Hours:        hours,
		Description:  description,
		HourlyRate:   hourlyRate,
		ProjectID:    projectObj.ID,
		ConsultantID: consultantObj.ID,
	}

	if err := repo.CreateTimeEntry(entry); err != nil {
		return fmt.Errorf("kunde inte spara tidsregistrering: %w", err)
	}

	cost := hours * hourlyRate
	fmt.Printf("✓ Tidsregistrering sparad!\n")
	fmt.Printf("  Datum: %s\n", entryDate.Format("2006-01-02"))
	fmt.Printf("  Konsult: %s (%.2f kr/h)\n", consultantObj.Name, hourlyRate)
	fmt.Printf("  Timmar: %.2f\n", hours)
	fmt.Printf("  Kostnad: %.2f kr\n", cost)
	fmt.Printf("  Projekt: %s\n", project)
	fmt.Printf("  Kund: %s\n", client)
	fmt.Printf("  Beskrivning: %s\n", description)

	return nil
}

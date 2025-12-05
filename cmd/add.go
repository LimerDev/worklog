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
	Short: "Add a new time entry",
	Long:  `Register a new time entry with hours, description, project and customer.`,
	RunE:  runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().Float64VarP(&hours, "hours", "t", 0, "Number of hours (required)")
	addCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the work (required)")
	addCmd.Flags().StringVarP(&project, "project", "p", "", "Project name (uses default if not specified)")
	addCmd.Flags().StringVarP(&client, "client", "c", "", "Customer name (uses default if not specified)")
	addCmd.Flags().StringVarP(&consultant, "consultant", "n", "", "Consultant name (uses default if not specified)")
	addCmd.Flags().Float64VarP(&hourlyRate, "rate", "r", 0, "Hourly rate (uses default if not specified)")
	addCmd.Flags().StringVarP(&date, "date", "D", "", "Date (YYYY-MM-DD, default: today)")

	addCmd.MarkFlagRequired("hours")
	addCmd.MarkFlagRequired("description")
}

func runAdd(cmd *cobra.Command, args []string) error {
	// Get configuration
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to read configuration: %w", err)
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
		return fmt.Errorf("consultant required (-n CONSULTANT or `worklog config set -n CONSULTANT`)")
	}
	if client == "" {
		return fmt.Errorf("customer required (-c CUSTOMER or `worklog config set -c CUSTOMER`)")
	}
	if project == "" {
		return fmt.Errorf("project required (-p PROJECT or `worklog config set -p PROJECT`)")
	}
	if hourlyRate <= 0 {
		return fmt.Errorf("hourly rate required (-r RATE or `worklog config set -r RATE`)")
	}

	var entryDate time.Time

	if date == "" {
		entryDate = time.Now()
	} else {
		entryDate, err = time.Parse("2006-01-02", date)
		if err != nil {
			return fmt.Errorf("invalid date format, use YYYY-MM-DD: %w", err)
		}
	}

	if hours <= 0 {
		return fmt.Errorf("hours must be greater than 0")
	}

	repo := database.NewRepository()

	// Get or create consultant
	consultantObj, err := repo.GetOrCreateConsultant(consultant)
	if err != nil {
		return fmt.Errorf("failed to get/create consultant: %w", err)
	}

	// Get or create customer
	customerObj, err := repo.GetOrCreateCustomer(client)
	if err != nil {
		return fmt.Errorf("failed to get/create customer: %w", err)
	}

	// Get or create project for this customer
	projectObj, err := repo.GetOrCreateProject(project, customerObj.ID)
	if err != nil {
		return fmt.Errorf("failed to get/create project: %w", err)
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
		return fmt.Errorf("failed to save time entry: %w", err)
	}

	cost := hours * hourlyRate
	fmt.Printf("✓ Time entry saved!\n")
	fmt.Printf("  Date: %s\n", entryDate.Format("2006-01-02"))
	fmt.Printf("  Consultant: %s (%.2f kr/h)\n", consultantObj.Name, hourlyRate)
	fmt.Printf("  Hours: %.2f\n", hours)
	fmt.Printf("  Cost: %.2f kr\n", cost)
	fmt.Printf("  Project: %s\n", project)
	fmt.Printf("  Customer: %s\n", client)
	fmt.Printf("  Description: %s\n", description)

	return nil
}

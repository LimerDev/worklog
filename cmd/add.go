package cmd

import (
	"fmt"
	"time"

	"github.com/LimerDev/worklog/internal/config"
	"github.com/LimerDev/worklog/internal/database"
	"github.com/LimerDev/worklog/internal/i18n"
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
	Short: "", // Set after i18n initialization
	Long:  "", // Set after i18n initialization
	RunE:  runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().Float64VarP(&hours, "hours", "t", 0, "")
	addCmd.Flags().StringVarP(&description, "description", "d", "", "")
	addCmd.Flags().StringVarP(&project, "project", "p", "", "")
	addCmd.Flags().StringVarP(&client, "client", "c", "", "")
	addCmd.Flags().StringVarP(&consultant, "consultant", "n", "", "")
	addCmd.Flags().Float64VarP(&hourlyRate, "rate", "r", 0, "")
	addCmd.Flags().StringVarP(&date, "date", "D", "", "")

	addCmd.MarkFlagRequired("hours")
	addCmd.MarkFlagRequired("description")
}

func localizeAddCommand() {
	addCmd.Short = i18n.T(i18n.KeyAddShort)
	addCmd.Long = i18n.T(i18n.KeyAddLong)

	addCmd.Flags().Lookup("hours").Usage = i18n.T(i18n.KeyAddFlagHours)
	addCmd.Flags().Lookup("description").Usage = i18n.T(i18n.KeyAddFlagDescription)
	addCmd.Flags().Lookup("project").Usage = i18n.T(i18n.KeyAddFlagProject)
	addCmd.Flags().Lookup("client").Usage = i18n.T(i18n.KeyAddFlagClient)
	addCmd.Flags().Lookup("consultant").Usage = i18n.T(i18n.KeyAddFlagConsultant)
	addCmd.Flags().Lookup("rate").Usage = i18n.T(i18n.KeyAddFlagRate)
	addCmd.Flags().Lookup("date").Usage = i18n.T(i18n.KeyAddFlagDate)
}

func runAdd(cmd *cobra.Command, args []string) error {
	// Get configuration
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrReadConfig), err)
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
		return fmt.Errorf("%s", i18n.T(i18n.KeyErrConsultantRequired))
	}
	if client == "" {
		return fmt.Errorf("%s", i18n.T(i18n.KeyErrCustomerRequired))
	}
	if project == "" {
		return fmt.Errorf("%s", i18n.T(i18n.KeyErrProjectRequired))
	}
	if hourlyRate <= 0 {
		return fmt.Errorf("%s", i18n.T(i18n.KeyErrRateRequired))
	}

	var entryDate time.Time

	if date == "" {
		now := time.Now()
		entryDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	} else {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrInvalidDateFormat), err)
		}
		entryDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
	}

	if hours <= 0 {
		return fmt.Errorf("%s", i18n.T(i18n.KeyErrHoursMustBePositive))
	}

	repo := database.NewRepository()

	// Get or create consultant
	consultantObj, err := repo.GetOrCreateConsultant(consultant)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrGetCreateConsultant), err)
	}

	// Get or create customer
	customerObj, err := repo.GetOrCreateCustomer(client)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrGetCreateCustomer), err)
	}

	// Get or create project for this customer
	projectObj, err := repo.GetOrCreateProject(project, customerObj.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrGetCreateProject), err)
	}

	// Check if there's an existing entry that matches all fields except hours
	existingEntry, err := repo.FindMatchingTimeEntry(entryDate, consultantObj.ID, projectObj.ID, description, hourlyRate)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrCheckExistingEntry), err)
	}

	var totalHours float64
	if existingEntry != nil {
		// Entry exists, update it by adding the hours
		if err := repo.UpdateTimeEntryHours(existingEntry.ID, hours); err != nil {
			return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrUpdateWorkLog), err)
		}
		totalHours = existingEntry.Hours + hours
		fmt.Println(i18n.T(i18n.KeyAddSuccessMerged))
	} else {
		// No matching entry, create a new one
		entry := &models.TimeEntry{
			Date:         entryDate,
			Hours:        hours,
			Description:  description,
			HourlyRate:   hourlyRate,
			ProjectID:    projectObj.ID,
			ConsultantID: consultantObj.ID,
		}

		if err := repo.CreateTimeEntry(entry); err != nil {
			return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrSaveWorkLog), err)
		}
		totalHours = hours
		fmt.Println(i18n.T(i18n.KeyAddSuccess))
	}

	cost := totalHours * hourlyRate
	fmt.Printf(i18n.T(i18n.KeyAddOutputDate)+"\n", entryDate.Format("2006-01-02"))
	fmt.Printf(i18n.T(i18n.KeyAddOutputConsultant)+"\n", consultantObj.Name)
	fmt.Printf(i18n.T(i18n.KeyAddOutputHours)+"\n", totalHours)
	fmt.Printf(i18n.T(i18n.KeyAddOutputRate)+"\n", hourlyRate)
	fmt.Printf(i18n.T(i18n.KeyAddOutputCost)+"\n", cost)
	fmt.Printf(i18n.T(i18n.KeyAddOutputProject)+"\n", project)
	fmt.Printf(i18n.T(i18n.KeyAddOutputCustomer)+"\n", client)
	fmt.Printf(i18n.T(i18n.KeyAddOutputDescription)+"\n", description)

	return nil
}

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/LimerDev/worklog/internal/database"
	"github.com/LimerDev/worklog/internal/i18n"
	"github.com/LimerDev/worklog/internal/output"
	"github.com/spf13/cobra"
)

var (
	getConsultant string
	getProject    string
	getCustomer   string
	getMonth      int
	getFromDate   string
	getToDate     string
	getDate       string
	getToday      bool
	getWeek       int
	getYear       int
	getOutput     string
	getOutputFile string
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Long:  "",
	RunE:  runGet,
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&getConsultant, "consultant", "n", "", "")
	getCmd.Flags().StringVarP(&getProject, "project", "p", "", "")
	getCmd.Flags().StringVarP(&getCustomer, "customer", "c", "", "")
	getCmd.Flags().IntVarP(&getMonth, "month", "m", 0, "")
	getCmd.Flags().StringVar(&getFromDate, "from", "", "")
	getCmd.Flags().StringVar(&getToDate, "to", "", "")
	getCmd.Flags().StringVarP(&getDate, "date", "D", "", "")
	getCmd.Flags().BoolVar(&getToday, "today", false, "")
	getCmd.Flags().IntVarP(&getWeek, "week", "w", 0, "")
	getCmd.Flags().IntVarP(&getYear, "year", "y", 0, "")
	getCmd.Flags().StringVarP(&getOutput, "output", "o", "table", "")
	getCmd.Flags().StringVar(&getOutputFile, "output-file", "", "")
}

func runGet(cmd *cobra.Command, args []string) error {
	// Handle --today flag
	if getToday {
		getDate = time.Now().Format("2006-01-02")
	}

	// Parse date filters
	var startDate, endDate time.Time

	// Handle week filter
	if getWeek > 0 {
		if getWeek < 1 || getWeek > 53 {
			return fmt.Errorf("%s", i18n.T(i18n.KeyErrWeekRange))
		}

		year := getYear
		if year == 0 {
			year = time.Now().Year()
		}

		// Find the first day of the week
		// Start from January 1st of the year and find the first Monday
		jan1 := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)

		// Calculate the Monday of week 1 (ISO 8601: week 1 is the first week with Thursday)
		// Find the first Thursday
		daysUntilThursday := (11 - int(jan1.Weekday())) % 7
		firstThursday := jan1.AddDate(0, 0, daysUntilThursday)

		// The Monday of week 1 is 3 days before the first Thursday
		firstMonday := firstThursday.AddDate(0, 0, -3)

		// Calculate the start date of the requested week
		startDate = firstMonday.AddDate(0, 0, 7*(getWeek-1))
		endDate = startDate.AddDate(0, 0, 7)
	} else if getDate != "" {
		parsedDate, err := time.Parse("2006-01-02", getDate)
		if err != nil {
			return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrInvalidDateFormat), err)
		}
		startDate = parsedDate
		endDate = parsedDate.AddDate(0, 0, 1)
	} else if getMonth > 0 {
		// Handle month filter
		if getMonth < 1 || getMonth > 12 {
			return fmt.Errorf("%s", i18n.T(i18n.KeyErrMonthRange))
		}

		year := getYear
		if year == 0 {
			year = time.Now().Year()
		}

		startDate = time.Date(year, time.Month(getMonth), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 1, 0)
	} else if getYear > 0 {
		// Handle year filter (when specified without month/week)
		startDate = time.Date(getYear, time.January, 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(1, 0, 0) // Next year
	} else {
		// Handle from/to date range
		if getFromDate != "" {
			parsedDate, err := time.Parse("2006-01-02", getFromDate)
			if err != nil {
				return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrInvalidDateFormat), err)
			}
			startDate = parsedDate
		}
		if getToDate != "" {
			parsedDate, err := time.Parse("2006-01-02", getToDate)
			if err != nil {
				return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrInvalidDateFormat), err)
			}
			endDate = parsedDate.AddDate(0, 0, 1) // Include the entire day
		}

		// If no date filters specified at all, default to current year and month
		if startDate.IsZero() && endDate.IsZero() {
			now := time.Now()
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
			endDate = startDate.AddDate(0, 1, 0)
		}
	}

	// Fetch entries
	repo := database.NewRepository()
	entries, err := repo.GetTimeEntriesByFilters(getConsultant, getProject, getCustomer, startDate, endDate)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrFetchWorkLogs), err)
	}

	if len(entries) == 0 {
		fmt.Println(i18n.T(i18n.KeyGetNoResults))
		return nil
	}

	// Determine output writer
	writer := os.Stdout
	if getOutputFile != "" {
		f, err := os.Create(getOutputFile)
		if err != nil {
			return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrCreateOutputFile), err)
		}
		defer f.Close()
		writer = f
	}

	// Get appropriate formatter
	format := output.Format(getOutput)
	formatter := output.GetFormatter(format)

	// Format and output results
	if err := formatter.Format(entries, writer); err != nil {
		return err
	}

	// Print success message if writing to file
	if getOutputFile != "" {
		fmt.Printf(i18n.T(i18n.KeyExportSuccess)+"\n", len(entries), getOutputFile)
	}

	return nil
}

func localizeGetCommand() {
	getCmd.Short = i18n.T(i18n.KeyGetShort)
	getCmd.Long = i18n.T(i18n.KeyGetLong)

	getCmd.Flags().Lookup("consultant").Usage = i18n.T(i18n.KeyGetFlagConsultant)
	getCmd.Flags().Lookup("project").Usage = i18n.T(i18n.KeyGetFlagProject)
	getCmd.Flags().Lookup("customer").Usage = i18n.T(i18n.KeyGetFlagCustomer)
	getCmd.Flags().Lookup("month").Usage = i18n.T(i18n.KeyGetFlagMonth)
	getCmd.Flags().Lookup("from").Usage = i18n.T(i18n.KeyGetFlagFromDate)
	getCmd.Flags().Lookup("to").Usage = i18n.T(i18n.KeyGetFlagToDate)
	getCmd.Flags().Lookup("date").Usage = i18n.T(i18n.KeyGetFlagDate)
	getCmd.Flags().Lookup("today").Usage = i18n.T(i18n.KeyGetFlagToday)
	getCmd.Flags().Lookup("week").Usage = i18n.T(i18n.KeyGetFlagWeek)
	getCmd.Flags().Lookup("year").Usage = i18n.T(i18n.KeyGetFlagYear)
	getCmd.Flags().Lookup("output").Usage = i18n.T(i18n.KeyGetFlagOutput)
	getCmd.Flags().Lookup("output-file").Usage = i18n.T(i18n.KeyGetFlagOutputFile)
}

package cmd

import (
	"fmt"
	"time"

	"github.com/LimerDev/worklog/internal/database"
	"github.com/LimerDev/worklog/internal/i18n"
	"github.com/LimerDev/worklog/internal/models"
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
			return fmt.Errorf(i18n.T(i18n.KeyErrWeekRange))
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
			return fmt.Errorf(i18n.T(i18n.KeyErrMonthRange))
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

	// Display results
	displayTimeEntries(entries)

	return nil
}

func displayTimeEntries(entries []models.TimeEntry) {
	type entrySummary struct {
		hours float64
		cost  float64
	}

	// Calculate column widths dynamically based on translated headers
	dateWidth := len(i18n.T(i18n.KeyGetHeaderDate))
	consultantWidth := len(i18n.T(i18n.KeyGetHeaderConsultant))
	projectWidth := len(i18n.T(i18n.KeyGetHeaderProject))
	customerWidth := len(i18n.T(i18n.KeyGetHeaderCustomer))
	descriptionWidth := len(i18n.T(i18n.KeyGetHeaderDescription))
	hoursWidth := len(i18n.T(i18n.KeyGetHeaderHours))
	rateWidth := len(i18n.T(i18n.KeyGetHeaderRate))
	costWidth := len(i18n.T(i18n.KeyGetHeaderCost))

	const (
		maxDateWidth        = 10
		maxConsultantWidth  = 20
		maxProjectWidth     = 25
		maxCustomerWidth    = 20
		maxDescriptionWidth = 40
	)

	for _, entry := range entries {
		// Text columns
		if len(entry.Date.Format("2006-01-02")) > dateWidth {
			dateWidth = len(entry.Date.Format("2006-01-02"))
		}
		if len(entry.Consultant.Name) > consultantWidth {
			consultantWidth = len(entry.Consultant.Name)
		}
		if len(entry.Project.Name) > projectWidth {
			projectWidth = len(entry.Project.Name)
		}
		if len(entry.Project.Customer.Name) > customerWidth {
			customerWidth = len(entry.Project.Customer.Name)
		}
		if len(entry.Description) > descriptionWidth {
			descriptionWidth = len(entry.Description)
		}

		// Calculate width for numeric values
		hoursStr := fmt.Sprintf("%.2f", entry.Hours)
		if len(hoursStr) > hoursWidth {
			hoursWidth = len(hoursStr)
		}
		rateStr := fmt.Sprintf("%.2f", entry.HourlyRate)
		if len(rateStr) > rateWidth {
			rateWidth = len(rateStr)
		}
		costStr := fmt.Sprintf("%.2f", entry.Hours*entry.HourlyRate)
		if len(costStr) > costWidth {
			costWidth = len(costStr)
		}
	}

	// Cap at max widths
	if dateWidth > maxDateWidth {
		dateWidth = maxDateWidth
	}
	if consultantWidth > maxConsultantWidth {
		consultantWidth = maxConsultantWidth
	}
	if projectWidth > maxProjectWidth {
		projectWidth = maxProjectWidth
	}
	if customerWidth > maxCustomerWidth {
		customerWidth = maxCustomerWidth
	}
	if descriptionWidth > maxDescriptionWidth {
		descriptionWidth = maxDescriptionWidth
	}

	consultantStats := make(map[string]*entrySummary)
	projectStats := make(map[string]*entrySummary)
	customerStats := make(map[string]*entrySummary)
	var totalHours float64
	var totalCost float64

	// Print header
	fmt.Printf("%-*s   %-*s   %-*s   %-*s   %-*s   %-*s   %-*s   %-*s\n",
		dateWidth, i18n.T(i18n.KeyGetHeaderDate),
		consultantWidth, i18n.T(i18n.KeyGetHeaderConsultant),
		hoursWidth, i18n.T(i18n.KeyGetHeaderHours),
		rateWidth, i18n.T(i18n.KeyGetHeaderRate),
		costWidth, i18n.T(i18n.KeyGetHeaderCost),
		projectWidth, i18n.T(i18n.KeyGetHeaderProject),
		customerWidth, i18n.T(i18n.KeyGetHeaderCustomer),
		descriptionWidth, i18n.T(i18n.KeyGetHeaderDescription))

	for _, entry := range entries {
		projectName := entry.Project.Name
		customerName := entry.Project.Customer.Name
		consultantName := entry.Consultant.Name
		hourlyRate := entry.HourlyRate
		cost := entry.Hours * hourlyRate

		fmt.Printf("%-*s   %-*s   %-*.2f   %-*.2f   %-*.2f   %-*s   %-*s   %-*s\n",
			dateWidth, entry.Date.Format("2006-01-02"),
			consultantWidth, truncate(consultantName, consultantWidth),
			hoursWidth, entry.Hours,
			rateWidth, hourlyRate,
			costWidth, cost,
			projectWidth, truncate(projectName, projectWidth),
			customerWidth, truncate(customerName, customerWidth),
			descriptionWidth, truncate(entry.Description, descriptionWidth))

		// Aggregate statistics
		projectName = entry.Project.Name
		customerName = entry.Project.Customer.Name
		consultantName = entry.Consultant.Name
		totalHours += entry.Hours
		totalCost += cost

		if _, exists := consultantStats[consultantName]; !exists {
			consultantStats[consultantName] = &entrySummary{}
		}
		consultantStats[consultantName].hours += entry.Hours
		consultantStats[consultantName].cost += cost

		if _, exists := projectStats[projectName]; !exists {
			projectStats[projectName] = &entrySummary{}
		}
		projectStats[projectName].hours += entry.Hours
		projectStats[projectName].cost += cost

		if _, exists := customerStats[customerName]; !exists {
			customerStats[customerName] = &entrySummary{}
		}
		customerStats[customerName].hours += entry.Hours
		customerStats[customerName].cost += cost
	}

	fmt.Printf("\n%s: %.2f\n", i18n.T(i18n.KeyGetTotalHours), totalHours)
	fmt.Printf("%s: %.2f kr\n", i18n.T(i18n.KeyGetTotalCost), totalCost)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
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
}

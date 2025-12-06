package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/LimerDev/worklog/internal/database"
	"github.com/spf13/cobra"
)

var (
	exportConsultant string
	exportProject    string
	exportCustomer   string
	exportMonth      string
	exportFromDate   string
	exportToDate     string
	exportDate       string
	exportToday      bool
	exportOutput     string
	exportFormat     string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export work logs to file",
	Long:  `Export work logs with flexible filtering to CSV format.`,
	RunE:  runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportConsultant, "consultant", "n", "", "Filter by consultant name")
	exportCmd.Flags().StringVarP(&exportProject, "project", "p", "", "Filter by project name")
	exportCmd.Flags().StringVarP(&exportCustomer, "customer", "c", "", "Filter by customer name")
	exportCmd.Flags().StringVarP(&exportMonth, "month", "m", "", "Filter by month (YYYY-MM)")
	exportCmd.Flags().StringVar(&exportFromDate, "from", "", "Filter from date (YYYY-MM-DD)")
	exportCmd.Flags().StringVar(&exportToDate, "to", "", "Filter to date (YYYY-MM-DD)")
	exportCmd.Flags().StringVarP(&exportDate, "date", "D", "", "Filter by specific date (YYYY-MM-DD)")
	exportCmd.Flags().BoolVar(&exportToday, "today", false, "Filter by today's date")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file (default: stdout)")
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "csv", "Export format (csv)")
}

func runExport(cmd *cobra.Command, args []string) error {
	// Note: We don't use defaults automatically for export command
	// User must explicitly specify filters they want

	// Handle --today flag
	if exportToday {
		exportDate = time.Now().Format("2006-01-02")
	}

	// Parse date filters
	var startDate, endDate time.Time

	// Handle specific date
	if exportDate != "" {
		parsedDate, err := time.Parse("2006-01-02", exportDate)
		if err != nil {
			return fmt.Errorf("invalid date format for --date, use YYYY-MM-DD: %w", err)
		}
		startDate = parsedDate
		endDate = parsedDate.AddDate(0, 0, 1)
	} else if exportMonth != "" {
		// Handle month filter
		parsedDate, err := time.Parse("2006-01", exportMonth)
		if err != nil {
			return fmt.Errorf("invalid month format for --month, use YYYY-MM: %w", err)
		}
		startDate = parsedDate
		endDate = parsedDate.AddDate(0, 1, 0)
	} else {
		// Handle from/to date range
		if exportFromDate != "" {
			parsedDate, err := time.Parse("2006-01-02", exportFromDate)
			if err != nil {
				return fmt.Errorf("invalid date format for --from, use YYYY-MM-DD: %w", err)
			}
			startDate = parsedDate
		}
		if exportToDate != "" {
			parsedDate, err := time.Parse("2006-01-02", exportToDate)
			if err != nil {
				return fmt.Errorf("invalid date format for --to, use YYYY-MM-DD: %w", err)
			}
			endDate = parsedDate.AddDate(0, 0, 1) // Include the entire day
		}
	}

	// Fetch entries
	repo := database.NewRepository()
	entries, err := repo.GetTimeEntriesByFilters(exportConsultant, exportProject, exportCustomer, startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed to fetch work logs: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No work logs found matching the filters.")
		return nil
	}

	// Export to CSV
	var output *os.File
	if exportOutput == "" {
		output = os.Stdout
	} else {
		f, err := os.Create(exportOutput)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		output = f
	}

	writer := csv.NewWriter(output)
	defer writer.Flush()

	// Write header
	header := []string{"DATE", "CONSULTANT", "PROJECT", "CUSTOMER", "DESCRIPTION", "HOURS", "RATE", "COST"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	var totalHours float64
	var totalCost float64

	// Write data rows
	for _, entry := range entries {
		cost := entry.Hours * entry.HourlyRate
		row := []string{
			entry.Date.Format("2006-01-02"),
			entry.Consultant.Name,
			entry.Project.Name,
			entry.Project.Customer.Name,
			entry.Description,
			fmt.Sprintf("%.2f", entry.Hours),
			fmt.Sprintf("%.2f", entry.HourlyRate),
			fmt.Sprintf("%.2f", cost),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
		totalHours += entry.Hours
		totalCost += cost
	}

	// Write totals row
	totalsRow := []string{
		"",
		"",
		"",
		"",
		"TOTAL",
		fmt.Sprintf("%.2f", totalHours),
		"",
		fmt.Sprintf("%.2f", totalCost),
	}
	if err := writer.Write(totalsRow); err != nil {
		return fmt.Errorf("failed to write totals row: %w", err)
	}

	if exportOutput != "" {
		fmt.Printf("✓ Exported %d work logs to %s\n", len(entries), exportOutput)
	}

	return nil
}

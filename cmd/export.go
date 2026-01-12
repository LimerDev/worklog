package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/LimerDev/worklog/internal/database"
	"github.com/LimerDev/worklog/internal/i18n"
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
	Short: "",
	Long:  "",
	RunE:  runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportConsultant, "consultant", "n", "", "")
	exportCmd.Flags().StringVarP(&exportProject, "project", "p", "", "")
	exportCmd.Flags().StringVarP(&exportCustomer, "customer", "c", "", "")
	exportCmd.Flags().StringVarP(&exportMonth, "month", "m", "", "")
	exportCmd.Flags().StringVar(&exportFromDate, "from", "", "")
	exportCmd.Flags().StringVar(&exportToDate, "to", "", "")
	exportCmd.Flags().StringVarP(&exportDate, "date", "D", "", "")
	exportCmd.Flags().BoolVar(&exportToday, "today", false, "")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "")
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "csv", "")
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
			return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrInvalidDateFormat), err)
		}
		startDate = parsedDate
		endDate = parsedDate.AddDate(0, 0, 1)
	} else if exportMonth != "" {
		// Handle month filter
		parsedDate, err := time.Parse("2006-01", exportMonth)
		if err != nil {
			return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrInvalidDateFormat), err)
		}
		startDate = parsedDate
		endDate = parsedDate.AddDate(0, 1, 0)
	} else {
		// Handle from/to date range
		if exportFromDate != "" {
			parsedDate, err := time.Parse("2006-01-02", exportFromDate)
			if err != nil {
				return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrInvalidDateFormat), err)
			}
			startDate = parsedDate
		}
		if exportToDate != "" {
			parsedDate, err := time.Parse("2006-01-02", exportToDate)
			if err != nil {
				return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrInvalidDateFormat), err)
			}
			endDate = parsedDate.AddDate(0, 0, 1) // Include the entire day
		}
	}

	// Fetch entries
	repo := database.NewRepository()
	entries, err := repo.GetTimeEntriesByFilters(exportConsultant, exportProject, exportCustomer, startDate, endDate)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrFetchWorkLogs), err)
	}

	if len(entries) == 0 {
		fmt.Println(i18n.T(i18n.KeyGetNoResults))
		return nil
	}

	// Export to CSV
	var output *os.File
	if exportOutput == "" {
		output = os.Stdout
	} else {
		f, err := os.Create(exportOutput)
		if err != nil {
			return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrCreateOutputFile), err)
		}
		defer f.Close()
		output = f
	}

	writer := csv.NewWriter(output)
	defer writer.Flush()

	// Write header
	header := []string{
		i18n.T(i18n.KeyGetHeaderDate),
		i18n.T(i18n.KeyGetHeaderConsultant),
		i18n.T(i18n.KeyGetHeaderProject),
		i18n.T(i18n.KeyGetHeaderCustomer),
		i18n.T(i18n.KeyGetHeaderDescription),
		i18n.T(i18n.KeyGetHeaderHours),
		i18n.T(i18n.KeyGetHeaderRate),
		i18n.T(i18n.KeyGetHeaderCost),
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrWriteCSVHeader), err)
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
			return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrWriteCSVRow), err)
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
		i18n.T(i18n.KeyExportTotal),
		fmt.Sprintf("%.2f", totalHours),
		"",
		fmt.Sprintf("%.2f", totalCost),
	}
	if err := writer.Write(totalsRow); err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrWriteTotalsRow), err)
	}

	if exportOutput != "" {
		fmt.Printf(i18n.T(i18n.KeyExportSuccess)+"\n", len(entries), exportOutput)
	}

	return nil
}

func localizeExportCommand() {
	exportCmd.Short = i18n.T(i18n.KeyExportShort)
	exportCmd.Long = i18n.T(i18n.KeyExportLong)

	exportCmd.Flags().Lookup("consultant").Usage = i18n.T(i18n.KeyGetFlagConsultant)
	exportCmd.Flags().Lookup("project").Usage = i18n.T(i18n.KeyGetFlagProject)
	exportCmd.Flags().Lookup("customer").Usage = i18n.T(i18n.KeyGetFlagCustomer)
	exportCmd.Flags().Lookup("month").Usage = i18n.T(i18n.KeyGetFlagMonth)
	exportCmd.Flags().Lookup("from").Usage = i18n.T(i18n.KeyGetFlagFromDate)
	exportCmd.Flags().Lookup("to").Usage = i18n.T(i18n.KeyGetFlagToDate)
	exportCmd.Flags().Lookup("date").Usage = i18n.T(i18n.KeyGetFlagDate)
	exportCmd.Flags().Lookup("today").Usage = i18n.T(i18n.KeyGetFlagToday)
	exportCmd.Flags().Lookup("output").Usage = i18n.T(i18n.KeyExportFlagOutput)
	exportCmd.Flags().Lookup("format").Usage = i18n.T(i18n.KeyExportFlagFormat)
}

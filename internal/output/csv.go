package output

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/LimerDev/worklog/internal/i18n"
	"github.com/LimerDev/worklog/internal/models"
)

// CSVFormatter formats entries as CSV
type CSVFormatter struct{}

// Format writes entries to the writer in CSV format
func (f *CSVFormatter) Format(entries []models.TimeEntry, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

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
	if err := csvWriter.Write(header); err != nil {
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
		if err := csvWriter.Write(row); err != nil {
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
	if err := csvWriter.Write(totalsRow); err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrWriteTotalsRow), err)
	}

	return nil
}

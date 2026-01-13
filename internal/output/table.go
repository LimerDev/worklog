package output

import (
	"fmt"
	"io"

	"github.com/LimerDev/worklog/internal/i18n"
	"github.com/LimerDev/worklog/internal/models"
)

// TableFormatter formats entries as a table
type TableFormatter struct{}

type entrySummary struct {
	hours float64
	cost  float64
}

// Format writes entries to the writer in table format
func (f *TableFormatter) Format(entries []models.TimeEntry, writer io.Writer) error {
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
	fmt.Fprintf(writer, "%-*s   %-*s   %-*s   %-*s   %-*s   %-*s   %-*s   %-*s\n",
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

		fmt.Fprintf(writer, "%-*s   %-*s   %-*.2f   %-*.2f   %-*.2f   %-*s   %-*s   %-*s\n",
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

	fmt.Fprintf(writer, "\n%s: %.2f\n", i18n.T(i18n.KeyGetTotalHours), totalHours)
	fmt.Fprintf(writer, "%s: %.2f kr\n", i18n.T(i18n.KeyGetTotalCost), totalCost)

	return nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

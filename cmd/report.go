package cmd

import (
	"fmt"
	"time"

	"github.com/LimerDev/worklog/internal/database"
	"github.com/spf13/cobra"
)

var (
	reportMonth string
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate monthly report",
	Long:  `Show a summary of all time entries for a given month.`,
	RunE:  runReport,
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.Flags().StringVarP(&reportMonth, "month", "m", "", "Month (YYYY-MM, default: current month)")
}

func runReport(cmd *cobra.Command, args []string) error {
	var year int
	var month time.Month

	if reportMonth == "" {
		now := time.Now()
		year = now.Year()
		month = now.Month()
	} else {
		parsedDate, err := time.Parse("2006-01", reportMonth)
		if err != nil {
			return fmt.Errorf("invalid month format, use YYYY-MM: %w", err)
		}
		year = parsedDate.Year()
		month = parsedDate.Month()
	}

	repo := database.NewRepository()
	entries, err := repo.GetTimeEntriesByMonth(year, month)
	if err != nil {
		return fmt.Errorf("failed to fetch time entries: %w", err)
	}

	if len(entries) == 0 {
		fmt.Printf("No time entries found for %s %d\n", month, year)
		return nil
	}

	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  TIME REPORT - %s %d\n", month, year)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	type consultantCost struct {
		name  string
		hours float64
		cost  float64
	}

	projectHours := make(map[string]float64)
	projectCost := make(map[string]float64)
	clientHours := make(map[string]float64)
	clientCost := make(map[string]float64)
	consultantStats := make(map[string]*consultantCost)
	var totalHours float64
	var totalCost float64

	fmt.Printf("%-12s %-15s %-8s %-10s %-12s %-20s %-20s %s\n", "Date", "Consultant", "Hours", "Rate", "Cost", "Project", "Customer", "Description")
	fmt.Printf("%-12s %-15s %-8s %-10s %-12s %-20s %-20s %s\n", "────────", "─────────────", "─────", "────", "────", "───────", "────────", "───────────")

	for _, entry := range entries {
		projectName := entry.Project.Name
		customerName := entry.Project.Customer.Name
		consultantName := entry.Consultant.Name
		hourlyRate := entry.HourlyRate
		cost := entry.Hours * hourlyRate

		fmt.Printf("%-12s %-15s %-8.2f %-10.2f %-12.2f %-20s %-20s %s\n",
			entry.Date.Format("2006-01-02"),
			truncate(consultantName, 15),
			entry.Hours,
			hourlyRate,
			cost,
			truncate(projectName, 20),
			truncate(customerName, 20),
			truncate(entry.Description, 40))

		// Aggregate statistics
		projectHours[projectName] += entry.Hours
		projectCost[projectName] += cost
		clientHours[customerName] += entry.Hours
		clientCost[customerName] += cost
		totalHours += entry.Hours
		totalCost += cost

		// Track consultant stats
		if _, exists := consultantStats[consultantName]; !exists {
			consultantStats[consultantName] = &consultantCost{name: consultantName}
		}
		consultantStats[consultantName].hours += entry.Hours
		consultantStats[consultantName].cost += cost
	}

	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  SUMMARY\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	fmt.Printf("Total hours: %.2f\n", totalHours)
	fmt.Printf("Total cost: %.2f kr\n\n", totalCost)

	fmt.Printf("Per consultant:\n")
	for consultant, stats := range consultantStats {
		fmt.Printf("  %-30s %.2f hours | %.2f kr\n", consultant, stats.hours, stats.cost)
	}

	fmt.Printf("\nPer project:\n")
	for project, hours := range projectHours {
		fmt.Printf("  %-30s %.2f hours | %.2f kr\n", project, hours, projectCost[project])
	}

	fmt.Printf("\nPer customer:\n")
	for client, hours := range clientHours {
		fmt.Printf("  %-30s %.2f hours | %.2f kr\n", client, hours, clientCost[client])
	}

	fmt.Printf("\n")

	return nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

package cmd

import (
	"fmt"
	"time"

	"github.com/LimerDev/worklog/internal/db"
	"github.com/spf13/cobra"
)

var (
	reportMonth string
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generera månadsrapport",
	Long:  `Visa en sammanställning av alla tidsregistreringar för en given månad.`,
	RunE:  runReport,
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.Flags().StringVarP(&reportMonth, "month", "m", "", "Månad (YYYY-MM, standard: aktuell månad)")
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
			return fmt.Errorf("ogiltigt månadsformat, använd YYYY-MM: %w", err)
		}
		year = parsedDate.Year()
		month = parsedDate.Month()
	}

	repo := db.NewRepository()
	entries, err := repo.GetTimeEntriesByMonth(year, month)
	if err != nil {
		return fmt.Errorf("kunde inte hämta tidsregistreringar: %w", err)
	}

	if len(entries) == 0 {
		fmt.Printf("Inga tidsregistreringar hittades för %s %d\n", month, year)
		return nil
	}

	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  TIDSRAPPORT - %s %d\n", month, year)
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

	fmt.Printf("%-12s %-15s %-8s %-12s %-20s %-20s %s\n", "Datum", "Konsult", "Timmar", "Kostnad", "Projekt", "Kund", "Beskrivning")
	fmt.Printf("%-12s %-15s %-8s %-12s %-20s %-20s %s\n", "─────────", "───────────────", "──────", "──────────", "────────────────", "────────────────", "────────────")

	for _, entry := range entries {
		projectName := entry.Project.Name
		customerName := entry.Project.Customer.Name
		consultantName := entry.Consultant.Name
		cost := entry.Hours * entry.Consultant.HourlyRate

		fmt.Printf("%-12s %-15s %-8.2f %-12.2f %-20s %-20s %s\n",
			entry.Date.Format("2006-01-02"),
			truncate(consultantName, 15),
			entry.Hours,
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
	fmt.Printf("  SAMMANFATTNING\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	fmt.Printf("Totalt antal timmar: %.2f\n", totalHours)
	fmt.Printf("Total kostnad: %.2f kr\n\n", totalCost)

	fmt.Printf("Per konsult:\n")
	for consultant, stats := range consultantStats {
		fmt.Printf("  %-30s %.2f timmar | %.2f kr\n", consultant, stats.hours, stats.cost)
	}

	fmt.Printf("\nPer projekt:\n")
	for project, hours := range projectHours {
		fmt.Printf("  %-30s %.2f timmar | %.2f kr\n", project, hours, projectCost[project])
	}

	fmt.Printf("\nPer kund:\n")
	for client, hours := range clientHours {
		fmt.Printf("  %-30s %.2f timmar | %.2f kr\n", client, hours, clientCost[client])
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

package output

import (
	"encoding/json"
	"io"

	"github.com/LimerDev/worklog/internal/models"
)

// JSONFormatter formats entries as JSON
type JSONFormatter struct{}

// JSONEntry represents a time entry in JSON format
type JSONEntry struct {
	Date        string  `json:"date"`
	Consultant  string  `json:"consultant"`
	Project     string  `json:"project"`
	Customer    string  `json:"customer"`
	Description string  `json:"description"`
	Hours       float64 `json:"hours"`
	HourlyRate  float64 `json:"hourly_rate"`
	Cost        float64 `json:"cost"`
}

// JSONOutput represents the complete JSON output
type JSONOutput struct {
	Entries    []JSONEntry `json:"entries"`
	TotalHours float64     `json:"total_hours"`
	TotalCost  float64     `json:"total_cost"`
	Count      int         `json:"count"`
}

// Format writes entries to the writer in JSON format
func (f *JSONFormatter) Format(entries []models.TimeEntry, writer io.Writer) error {
	var jsonEntries []JSONEntry
	var totalHours float64
	var totalCost float64

	for _, entry := range entries {
		cost := entry.Hours * entry.HourlyRate
		jsonEntry := JSONEntry{
			Date:        entry.Date.Format("2006-01-02"),
			Consultant:  entry.Consultant.Name,
			Project:     entry.Project.Name,
			Customer:    entry.Project.Customer.Name,
			Description: entry.Description,
			Hours:       entry.Hours,
			HourlyRate:  entry.HourlyRate,
			Cost:        cost,
		}
		jsonEntries = append(jsonEntries, jsonEntry)
		totalHours += entry.Hours
		totalCost += cost
	}

	output := JSONOutput{
		Entries:    jsonEntries,
		TotalHours: totalHours,
		TotalCost:  totalCost,
		Count:      len(entries),
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

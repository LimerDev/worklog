package output

import (
	"io"

	"github.com/LimerDev/worklog/internal/models"
)

// Format represents the output format type
type Format string

const (
	FormatTable Format = "table"
	FormatCSV   Format = "csv"
	FormatJSON  Format = "json"
)

// Formatter is the interface for all output formatters
type Formatter interface {
	Format(entries []models.TimeEntry, writer io.Writer) error
}

// GetFormatter returns the appropriate formatter for the given format
func GetFormatter(format Format) Formatter {
	switch format {
	case FormatCSV:
		return &CSVFormatter{}
	case FormatJSON:
		return &JSONFormatter{}
	case FormatTable:
		fallthrough
	default:
		return &TableFormatter{}
	}
}

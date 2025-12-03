package models

import (
	"time"
)

// TimeEntry represents a time tracking entry for a project
type TimeEntry struct {
	ID           uint       `gorm:"primaryKey"`
	Date         time.Time  `gorm:"not null;index"`
	Hours        float64    `gorm:"not null"`
	Description  string     `gorm:"type:text;not null"`
	ProjectID    uint       `gorm:"not null;index"`
	Project      Project    `gorm:"foreignKey:ProjectID"`
	ConsultantID uint       `gorm:"not null;index"`
	Consultant   Consultant `gorm:"foreignKey:ConsultantID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

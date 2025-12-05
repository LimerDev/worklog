package models

import "time"

// Consultant represents a consultant/developer that can log time
type Consultant struct {
	ID          uint        `gorm:"primaryKey"`
	Name        string      `gorm:"uniqueIndex;not null"`
	Active      bool        `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	TimeEntries []TimeEntry `gorm:"foreignKey:ConsultantID"`
}

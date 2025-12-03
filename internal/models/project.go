package models

import "time"

// Project represents a project belonging to a customer
type Project struct {
	ID          uint        `gorm:"primaryKey"`
	Name        string      `gorm:"not null;index"`
	CustomerID  uint        `gorm:"not null;index"`
	Customer    Customer    `gorm:"foreignKey:CustomerID"`
	Active      bool        `gorm:"default:true"`
	Description string      `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	TimeEntries []TimeEntry `gorm:"foreignKey:ProjectID"`
}

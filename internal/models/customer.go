package models

import "time"

// Customer represents a client/customer
type Customer struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"uniqueIndex;not null"`
	Active    bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Projects  []Project `gorm:"foreignKey:CustomerID"`
}

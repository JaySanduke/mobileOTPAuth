package model

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	MobileNumber string    `gorm:"uniqueIndex" json:"mobile_number"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Fingerprint  string    `json:"fingerprint"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

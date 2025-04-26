package model

import (
	"time"
)

type UserLoginSession struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       string    `json:"user_id"`
	MobileNumber string    `gorm:"uniqueIndex" json:"mobile_number"`
	Fingerprint  string    `json:"fingerprint"`
	LastLogin    time.Time `json:"last_login"`
	AccessToken  string    `json:"access_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

package types

import (
	"time"

	"gorm.io/gorm"
)

type Contact struct {
	ID             int    `gorm:"primaryKey"`
	PhoneNumber    string `gorm:"index"`
	Email          string `gorm:"index"`
	LinkedID       int
	LinkPrecedence string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
}

type Request struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

type Response struct {
	PrimaryContactID    int      `json:"primaryContatctId"`
	Emails              []string `json:"emails"`
	PhoneNumbers        []string `json:"phoneNumbers"`
	SecondaryContactIDs []int    `json:"secondaryContactIds"`
}

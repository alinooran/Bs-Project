package models

import (
	"gorm.io/gorm"
	"time"
)

type Request struct {
	gorm.Model
	Date              time.Time `json:"date"`
	Sent              bool      `json:"sent"`
	FinalStatus       *bool     `json:"final_status"`
	UserID            uint
	UserIdForApproval *uint
	Workflows         []Workflow `json:"workflows"`
	Guests            []Guest    `json:"guests"`
}

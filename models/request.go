package models

import (
	"gorm.io/gorm"
	"time"
)

type Request struct {
	gorm.Model
	Date              time.Time `json:"date"`
	Sent              bool      `json:"sent"`
	DeanApproval      *bool     `json:"dean_approval"`
	SecurityApproval  *bool     `json:"security_approval"`
	FinalStatus       *bool     `json:"final_status"`
	UserID            uint      `json:"user_id"`
	UserIdForApproval *uint
	Workflows         []Workflow `json:"workflows"`
	Guests            []Guest    `json:"guests"`
}

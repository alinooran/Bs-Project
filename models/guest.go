package models

import (
	"gorm.io/gorm"
	"time"
)

type Guest struct {
	gorm.Model
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	StartTime    string     `json:"start_time"`
	EndTime      string     `json:"end_time"`
	NationalCode string     `json:"national_code"`
	IsEntered    bool       `json:"is_entered"`
	EnteredDate  *time.Time `json:"entered_date"`
	RequestID    uint
}

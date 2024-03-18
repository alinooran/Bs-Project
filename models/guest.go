package models

import (
	"gorm.io/gorm"
	"time"
)

type Guest struct {
	gorm.Model
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Date      time.Time `json:"date"`
	Phone     string    `json:"phone"`
	RequestID uint
}

package models

import "gorm.io/gorm"

type Workflow struct {
	gorm.Model
	SenderName  string `json:"sender_name"`
	Step        string `json:"step"`
	Status      *bool  `json:"status"`
	Description string `json:"description"`
	RequestID   uint
}

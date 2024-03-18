package models

import (
	"gorm.io/gorm"
	"time"
)

type Request struct {
	gorm.Model
	DepartmentApproval *bool `json:"department_approval"`
	SecurityApproval   *bool `json:"security_approval"`
	//FinalApproval      bool       `json:"final_approval"`
	Description       string     `json:"description"`
	Sent              bool       `json:"sent"`
	SentDate          *time.Time `json:"sent_date"`
	UserID            uint
	UserIdForApproval *uint
	Guests            []Guest `json:"guests"`
}

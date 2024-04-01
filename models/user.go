package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"unique"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email" gorm:"unique"`
	DepartmentID uint
	ForApproval  []Request `gorm:"foreignKey:UserIdForApproval"`
	MyRequests   []Request
}

package models

type Role struct {
	ID    uint   `gorm:"primarykey"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

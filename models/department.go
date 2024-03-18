package models

type Department struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Name  string `json:"name"`
	Users []User
}

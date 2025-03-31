package dto

import "gorm.io/gorm"

const (
	UsersPageSizeDefault = 20
	UsersPageDefault     = 1
)

type UsersListFilter struct {
	Page     int
	PageSize int
}

type User struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:255"`
	Email string `gorm:"unique"`
}

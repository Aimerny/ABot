package models

import "gorm.io/gorm"

type Priority int

const (
	P_HIGH   = 0
	P_MIDDLE = 1
	P_LOW    = 2
)

type Restaurant struct {
	gorm.Model
	Name     string
	Priority Priority
}

func (r *Restaurant) TableName() string {
	return "tb_rest"
}

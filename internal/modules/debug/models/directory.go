package models

import "time"

type Directory struct {
	Id        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null"`
	Remark    string `gorm:"size:255"`
	ParentId  uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

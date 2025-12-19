package models

import "time"

type Program struct {
	Id        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"Name" gorm:"size:255;not null"`
	Code      string `json:"code" gorm:"type:text"`
	Sort      uint   `json:"sort"`
	CreatedAt time.Time
	UpdatedAt time.Time
	PackageId uint `json:"packageid"`
}

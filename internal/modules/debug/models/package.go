package models

import "time"

type Package struct {
	Id          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"size:255"`
	Description string `json:"description" gorm:"type:text"`
	DirectoryId uint   `json:"directoryid"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

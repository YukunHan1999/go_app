package models

import "time"

type DebugInfo struct {
	Id           uint `json:"id" gorm:"primaryKey"`
	LineNo       int  `json:"lineno"`
	Attachmentid uint `json:"attachmentid" gorm:"size:255"`
	Sort         uint `json:"sort"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ProgramId    uint `json:"programid"`
}

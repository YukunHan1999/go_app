package models

import "time"

type EnvData struct {
	Id           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `gorm:"size:255;not null"`
	Remark       string    `gorm:"size:255"`
	Attachmentid uint      `json:"attachmentid" gorm:"size:255"`
	CreatedAt    time.Time `json:"createtime"`
	UpdatedAt    time.Time `json:"updatetime"`
	Pid          uint      `json:"pid"`
}

type EnvDataDTO struct {
	Id            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"size:255;not null"`
	Remark        string    `json:"desc" gorm:"size:255"`
	Attachmentid  uint      `json:"attachmentid" gorm:"size:255"`
	Attachmenturl string    `json:"attachmenturl" gorm:"size:255"`
	CreatedAt     time.Time `json:"createtime"`
	UpdatedAt     time.Time `json:"updatetime"`
	Pid           uint      `json:"dirid"`
}

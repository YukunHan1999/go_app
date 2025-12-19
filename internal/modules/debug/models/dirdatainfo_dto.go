package models

import (
	"time"
)

type RemovePkgOrDir struct {
	Id     uint `json:"id"`
	IsFile uint `json:"isdocx"`
}

type DirData struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	ParentId    uint   `json:"dirId"`
}

type DirDataInfo struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"desc"`
	CreatedAt   time.Time `json:"createdate"`
	UpdatedAt   time.Time `json:"updatedate"`
	ParentId    uint      `json:"dirId"`
	IsFile      uint8     `json:"isdocx"`
	Deleted     uint8     `json:"deleted"`
}

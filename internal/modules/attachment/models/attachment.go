package models

type Attachment struct {
	Id   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"size:255"`
	Type string `json:"type"`
	Url  string `json:"url" gorm:"size:255"`
}

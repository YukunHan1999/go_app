package models

type Codetemplate struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Lang string `json:"lang" gorm:"size:255;not null"`
	Name string `json:"name"`
	Code string `json:"code" gorm:"type:text"`
	Tips string `json:"tips"`
}

// TableName sets custom table name
func (Codetemplate) TableName() string {
	return "codetemplate" // 👈 map this struct to table "t_user"
}

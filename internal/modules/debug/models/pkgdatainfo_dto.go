package models

// everyline debuginfo
type DbgInfo struct {
	Id      uint   `json:"dbgid"`
	LineNo  int    `json:"lineno"`
	Sort    uint   `json:"dbgsort"`
	Attid   uint   `json:"attid"`
	AttName string `json:"attname"`
	AttType string `json:"atttype"`
	AttUrl  string `json:"atturl"`
	PgmId   uint
}

// program info
type PgmDataInfo struct {
	Id       uint      `json:"pgmid"`
	Name     string    `json:"pgmname"`
	Code     string    `json:"pgmcode"`
	Sort     uint      `json:"pgmsort"`
	DbgArray []DbgInfo `json:"dbg"`
	PkgId    uint
}

// package info
type PkgDataInfo struct {
	Id           uint          `json:"pkgid"`
	Name         string        `json:"pkgname"`
	Description  string        `json:"pkgdesc"`
	DirectoryId  uint          `json:"dirid"`
	ProgramArray []PgmDataInfo `json:"pgm"`
}

type AddPkgDto struct {
	Pkg    *PkgDataInfo `json:"pkgData"`
	AttIds []uint       `json:"attGarbage"`
}

type CleakPkgDto struct {
	Id     uint   `json:"pkgid"`
	AttIds []uint `json:"attGarbage"`
}

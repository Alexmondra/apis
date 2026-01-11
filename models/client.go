package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	Empresa  string `gorm:"type:varchar(250)" json:"empresa"`
	Ruc      string `gorm:"type:varchar(11);uniqueIndex" json:"ruc"`
	ApiKey   string `gorm:"type:varchar(100);uniqueIndex" json:"api_key"`
	Limit    int    `gorm:"default:100" json:"limit"`
	Usage    int    `gorm:"default:0" json:"usage"`
	Status   bool   `gorm:"default:true" json:"status"` // true = activo
}
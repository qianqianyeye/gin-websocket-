package model

import (
	"AdPushServer_Go/src/webgo"
)

type AndroidScreen struct {
	ID int `gorm:"column:id" json:"id"`
	AndroidModel string `gorm:"column:android_model" json:"android_model"`
	AndroidNumber string `gorm:"column:android_number" json:"android_number"`
	CreateTime webgo.JsonDateTime `gorm:"column:create_time" json:"create_time"`
	Status int `gorm:"column:status" json:"status"`
	IsBind int `gorm:"column:is_bind" json:"is_bind"`
	Version string `gorm:"column:version" json:"version"`
}

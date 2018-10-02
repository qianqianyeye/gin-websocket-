package model

import (
	"AdPushServer_Go/src/webgo"
)

type LogAndroidVersion struct {
	ID int64 `gorm:"column:id" json:"id"`
	AppVersion string `gorm:"column:app_version" json:"app_version"`
	Url string `gorm:"column:url" json:"url"`
	CreateTime webgo.JsonDateTime `gorm:"column:create_time" json:"create_time"`
	Remarks string `gorm:"column:remarks" json:"remarks"`
	Type int64 `gorm:"column:type" json:"type"`
	IsAllCheck int64 `gorm:"column:is_all_check" json:"is_all_check"`
}
package model

import (
	"AdPushServer_Go/src/webgo"
)

type Advert struct {
	ID int `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
	Model int64 `gorm:"column:model" json:"model"`
	ScrollTime int `gorm:"column:scroll_time" json:"scrollTime"`
	CreateTime webgo.JsonDateTime `gorm:"column:create_time" json:"createTime"`
	UpdateTime webgo.JsonDateTime `gorm:"column:update_time" json:"updateTime"`
	Status int `gorm:"column:status" json:"status"`
	AdvertContent  []AdvertContent    `gorm:"foreignkey:AdvertID" json:"advertContent"`
	AdvertDownload []AdvertDownload   `json:"advertDownload"`
	AndroidScreen  []AndroidScreen    `gorm:"many2many:android_screen json:"androidScreen";"`
}

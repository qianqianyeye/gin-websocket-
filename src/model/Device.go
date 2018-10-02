package model

import (
	"AdPushServer_Go/src/webgo"
)

type Device struct {
	ID           int                `json:"id"`
	HardwareID   string             `json:"hardwareId"`
	ModelID      uint               `json:"modelId"`
	IsBind       byte               `json:"isBind"`
	ReleaseDate  webgo.JsonDateTime `json:"releaseDate"`
	StoreID      uint               `json:"storeId"`
	Status       uint               `json:"status"`
	CreateTime   webgo.JsonDateTime `json:"createTime"`
	UpdateTime   webgo.JsonDateTime `json:"updateTime"`
	DeviceNumber string             `json:"deviceNumber"`
}

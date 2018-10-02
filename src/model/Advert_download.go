package model

import (
	"AdPushServer_Go/src/webgo"
)

type AdvertDownload struct {
	ID               int                `json:"id"`
	AdvertID         int                `json:"advertId"`
	AndroidScreenID  int                `json:"androidScreenId"`
	DownloadProgress int8               `json:"downloadProgress"`
	CreateTime       webgo.JsonDateTime `json:"createTime"`
	DeviceID         int                `json:"deviceId"`
	MerchantID       int                `json:"merchantId"`
	StoreID          int                `json:"storeId"`
	UpdateTime       webgo.JsonDateTime `json:"updateTime"`
}

package model

import (
	"AdPushServer_Go/src/webgo"
)

type DeviceModular struct {
	ID          uint               `json:"id"`
	DeviceID    uint               `json:"deviceId"`
	ModularID   int               `json:"modularId"`
	ModularType int               `json:"modularType"`
	InstallSite int8               `json:"installSite"`
	ConType     uint               `json:"conType"`
	CreateTime  webgo.JsonDateTime `json:"createTime"`
	Coin        int                `json:"coin"`
	CtrlFall    int                `json:"ctrlFall"`
	DefaultFall int                `json:"defaultFall"`
	ActualFall  int                `json:"actualFall"`
}

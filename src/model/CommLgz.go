package model

import (
	"AdPushServer_Go/src/webgo"
)

type CommLgz struct {
	ID           int                `json:"id"`
	LgzID        string             `json:"lgz_id"`
	IsBind       int8               `json:"is_bind"`
	Status       int8               `json:"status"`
	ErrCode      string             `json:"err_code"`
	CreateTime   webgo.JsonDateTime `json:"create_time"`
	UpdateTime   webgo.JsonDateTime `json:"update_time"`
	Time         int8               `json:"time"`
	Strong       int8               `json:"strong"`
	Weak         int8               `json:"weak"`
	Modular_type string             `json:"modular_type"`
	Install_site string             `json:"install_site"`
	Coin         string             `json:"coin"`

	DeviceModular DeviceModular
}

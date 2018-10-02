package model

import (
	"AdPushServer_Go/src/webgo"
)

type MerchantStore struct {
	ID         int                `json:"id"`
	MerchantID int                `json:"merchant_id"`
	StoreID    int                `json:"store_id"`
	CreateTime webgo.JsonDateTime `json:"create_time"`
}

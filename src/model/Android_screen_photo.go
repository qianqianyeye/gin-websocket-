package model

import "AdPushServer_Go/src/webgo"

type Android_screen_photo struct {
	Id                int            `json:"id"`
	Android_screen_id int            `json:"androidScreenId"`
	Img               string         `json:"img"`
	Type              int            `json:"type"` //1截屏，2照相
	CreatedAt         webgo.JsonTime `json:"createdAt"`
	UpdatedAt         webgo.JsonTime `json:"updatedAt"`
}

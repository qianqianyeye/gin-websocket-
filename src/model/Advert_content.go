package model

type AdvertContent struct {
	ID         int    `json:"id"`
	AdvertID   string `json:"advertId"`
	URL        string `json:"url"`
	AreaType   int8   `json:"areaType"`
	AreaWidth  int    `json:"areaWidth"`
	AreaHeight int    `json:"areaHeight"`
	X          int    `json:"x"`
	Y          int    `json:"y"`
	Advert     Advert `gorm:"foreignkey:AdvertID" json:"advert"`
	Flag       bool
}

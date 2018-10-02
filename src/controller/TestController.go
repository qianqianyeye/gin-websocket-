package controller

import (
	"AdPushServer_Go/src/db"
	"AdPushServer_Go/src/model"
	"AdPushServer_Go/src/webgo"
	"github.com/gin-gonic/gin"
)

type TestController struct {
	webgo.Controller
}

func (ctrl *TestController) Router(router *gin.Engine) {
	r := router.Group("test")
	//r.GET("advert",ctrl.sendAdById)
	r.POST("test", ctrl.test)

}

func (ctrl *TestController) test(c *gin.Context) {

	var android []model.AndroidScreen
	//rows,_:=db.SqlDB.Table("android_screen").Rows()
	db.SqlDB.Raw("SELECT * FROM android_screen").Select("*").Scan(&android)
	db.SqlDB.Table("android_screen").Related("advert_download", "android_screen_id")

}

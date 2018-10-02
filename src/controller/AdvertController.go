package controller

import (
	"AdPushServer_Go/src/db"
	"AdPushServer_Go/src/model"
	"AdPushServer_Go/src/service"
	"AdPushServer_Go/src/socketio"
	"AdPushServer_Go/src/webgo"
	"github.com/gin-gonic/gin"
	"net/http"
	"AdPushServer_Go/src/middleware"
	"fmt"
)

var adverService webservice.AdverService

type AdvertController struct {
	webgo.Controller
}

func (ctrl *AdvertController) Router(router *gin.Engine) {
	r := router.Group("api/v1", middleware.ClawMiddle)
	//r.GET("advert",ctrl.sendAdById)
	r.POST("advert", ctrl.sendAdById)
	r.POST("screen/screenShot", ctrl.screenShot)
	r.POST("screen/cameraMonitor", ctrl.cameraMonitor)
	r.POST("screen/qrCode", ctrl.qrCode)
}

func (ctrl *AdvertController) sendAdById(ctx *gin.Context) {
	s := goGetAdvert(ctx)
	temp := []int{}
	for _, val := range s.AdvertDownload {
		//验证是否要发送广告
		flag := adverService.Validata(val)
		if flag {
			temp = append(temp, val.AndroidScreenID)
		}
	}
	if len(temp) > 0 {
		socket := adverService.GetSocketIdKey(temp) //拿到AndroidNumber
		//redis :=db.CfRedis[0]
		if len(socket) > 0 {
			socketio.ParseAdJson(&s, socket, "updateAdModel")
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": "ok",
		"status":0,
	})
}

func goGetAdvert(ctx *gin.Context) model.Advert {
	//id := ctx.PostForm("id")
	var s model.Advert
	if _, ok := ctx.Keys["id"]; ok {
		id := webgo.GetResult(ctx.Keys["id"])
		s = adverService.GetAdById(id)
		return s
	}
	return s
}

func (ctrl *AdvertController) screenShot(ctx *gin.Context) {
	//获取传送的数组
	var arr []string
	if _, ok := ctx.Keys["android_screens"]; ok {
		android_screens := webgo.GetArr(ctx.Keys["android_screens"])
		redis := db.CfRedis[0]
		for _, val := range android_screens {
			androidNumber := redis.HGet("device:number", val).Val()
			clientId := redis.Get(androidNumber).Val()
			if clientId == "" {
				arr = append(arr, androidNumber)
				continue
			}
			sign := socketio.GetSign(androidNumber, "screenShot")
			result :=make(map[string]interface{})
			result["clientId"]=clientId
			result["data"]=gin.H{
				"socketType": "screenShot",
				"deviceID":   androidNumber,
				"data":       gin.H{},
				"sign":       sign,
				"status":1,
			}
			socketio.AndroidServer.Msg(result)
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status": 0,
			"data":   arr,
		})
	}

}
func (ctrl *AdvertController) cameraMonitor(ctx *gin.Context) {
	if _, ok := ctx.Keys["android_screens"]; ok {
		android_screens := webgo.GetArr(ctx.Keys["android_screens"])
		//获取传送的数组
		//var arr []string
		arr :=[]string{}
		redis := db.CfRedis[0]
		for _, val := range android_screens {
			androidNumber := redis.HGet("device:number", val).Val()
			clientId := redis.Get(androidNumber).Val()
			if clientId == "" {
				arr = append(arr, androidNumber)
				continue
			}
			sign := socketio.GetSign(androidNumber, "cameraMonitor")
			result :=make(map[string]interface{})
			result["clientId"]=clientId
			result["data"]=gin.H{
				"socketType": "cameraMonitor",
				"deviceID":   androidNumber,
				"data":       gin.H{},
				"sign":       sign,
				"status":1,
			}
			socketio.AndroidServer.Msg(result)
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status": 0,
			"data":   arr,
		})
	}
}

func (ctrl *AdvertController) qrCode(ctx *gin.Context) {
	if _, ok := ctx.Keys["id"]; ok {
		id := webgo.GetResult(ctx.Keys["id"])
		redis := db.CfRedis[0]
		androidNumber := redis.HGet("device:number", id).Val()
		clientId := redis.Get(androidNumber).Val()
		sign := socketio.GetSign(androidNumber, "deviceQrCode")
		qrcAdd := model.Claw + "?device_number=" + id
		for _,val := range model.QrcAndroidNum {
			fmt.Println(val)
			if id == val {
				qrcAdd = model.Claw + "miniapp?device_number=" + id
			}
		}
		result :=make(map[string]interface{})
		result["clientId"]=clientId
		result["data"]=gin.H{
			"socketType": "deviceQrCode",
			"deviceID":   androidNumber,
			"data":       qrcAdd,
			"sign":       sign,
			"status":1,
		}
		socketio.AndroidServer.Msg(result)
	}
}

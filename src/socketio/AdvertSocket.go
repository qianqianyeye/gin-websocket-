package socketio

import (
	"AdPushServer_Go/src/db"
	"AdPushServer_Go/src/model"
	"AdPushServer_Go/src/service"
	"AdPushServer_Go/src/webgo"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type AdvertSocket struct {
}

var deviceService webservice.DeviceService
var advertService webservice.AdverService

func AppVersion(dat map[string]interface{}, c *Client, sign string) {
	data := dat["data"].(map[string]interface{})
	version :=webgo.GetResult(data["version"])
	deviceId := webgo.GetResult(dat["deviceID"])
	redis := db.CfRedis[0]
	androidScreenId :=redis.HGet("device:id",deviceId).Val()
	var logAndroidVersion model.LogAndroidVersion
	db.SqlDB.Table("log_android_version").Select("*").Order(" create_time desc").Limit(1).Scan(&logAndroidVersion)
	if version==logAndroidVersion.AppVersion {
		db.SqlDB.Exec("update android_screen set version = ? where id = ?", logAndroidVersion.AppVersion,androidScreenId )
	}
		c.Write(gin.H{
			"sign":       sign,
			"socketType": "appVersion",
			"deviceID":   deviceId,
			"status":     1, //处理结果1为成功，0为失败
			"data": gin.H{
				"version":logAndroidVersion.AppVersion ,
				"updateType":logAndroidVersion.Type, //0乐关注安卓模块，1普通广告模块
				"appUrl":logAndroidVersion.Url,
			},
		})
}

//创建AndroidNumber
func CreateDeviceID(dat map[string]interface{}, c *Client, sign string) {
	data := dat["data"].(map[string]interface{})
	androidID := webgo.GetResult(data["androidID"])
	androidScreen := deviceService.CreateDevice(androidID)
	SetIdToNumber()
	c.Write(gin.H{
		"sign":       sign,
		"socketType": "createDeviceID",
		"deviceID":   androidScreen.AndroidNumber,
		"status":     1, //处理结果1为成功，0为失败
		"data": gin.H{
			"androidID": androidID,
		},
	})
}

//心跳
func HeartImpulse(dat map[string]interface{}, c *Client, sign string) {
	deviceId := webgo.GetResult(dat["deviceID"])
	//fmt.Println(deviceId)
	online(deviceId)
	setRedisKey(deviceId, c.Id) //保存socketID与AndroidNumber的关系
	setRedisKey(c.Id, deviceId)
	c.Write(gin.H{
		"status":     1, //处理结果1为成功，0为失败
		"socketType": "heartImpulse",
		"deviceID":   deviceId,
	})
}

func UpdateModel(dat map[string]interface{}, c *Client, sign string) {

	deviceId := webgo.GetResult(dat["deviceID"])
	data := dat["data"].(map[string]interface{})
	adUpdateID := webgo.GetResult(data["adUpdateID"])
	adDownloadProgress := webgo.GetResult(data["adDownloadProgress"])
	redis := db.CfRedis[0]
	android_screen_id := redis.HGet("device:id", deviceId).Val()
	//定时器处理存取的进度
	updateData := adUpdateID + "," + adDownloadProgress + "," + android_screen_id
	//存入数组到redis
	redis.RPush("updateAdModel", updateData)
}

//广告更新对比
func CheckUpdateID(dat map[string]interface{}, c *Client, sign string) {
	deviceId := webgo.GetResult(dat["deviceID"])
	data := dat["data"].(map[string]interface{})
	adUpdateID := webgo.GetResult(data["adUpdateID"])
	redis := db.CfRedis[0]
	android_screen_id := redis.HGet("device:id", deviceId).Val()
	//获取最新的一条广告下载进度
	adverDownLoad := advertService.FindAdnroidScreenOne(android_screen_id)
	advertId := adverDownLoad.AdvertID
	//数据库中没有匹配的广告
	if advertId == 0 {
		c.Write(gin.H{
			"sign":       sign,
			"status":     0,
			"socketType": "checkUpdateID",
			"message":    "你使用的是默认广告！",
			"deviceID":   deviceId,
			"data": gin.H{
				"adUpdateID": adUpdateID,
			},
		})
	}
	//广告Id跟adUpdateId是否相等 相等为最新广告
	if strconv.Itoa(advertId) == adUpdateID {
		c.Write(gin.H{
			"sign":       sign,
			"status":     0,
			"socketType": "checkUpdateID",
			"deviceID":   deviceId,
			"data": gin.H{
				"adUpdateID": adUpdateID,
			},
		})
	} else {
		advert := advertService.GetAdById(strconv.Itoa(advertId))
		var key []string
		key = append(key, deviceId)
		//下发新广告
		ParseAdJson(&advert, key, "checkUpdateID")
	}
}

//截屏
func ScreenShot(dat map[string]interface{}, c *Client, sign string) {
	deviceID := webgo.GetResult(dat["deviceID"])
	data := dat["data"].(map[string]interface{})
	imageUrl := webgo.GetResult(data["imageUrl"])
	redis := db.CfRedis[0]
	android_screen_id := redis.HGet("device:id", deviceID).Val()
	var androidScreenPhoto model.Android_screen_photo
	asId, _ := strconv.Atoi(android_screen_id)
	androidScreenPhoto.Android_screen_id = asId
	androidScreenPhoto.Img = imageUrl
	androidScreenPhoto.Type = 1
	advertService.InsertScreenPhoto(androidScreenPhoto)
	c.Write(gin.H{
		"sign":    sign,
		"status":  1, //处理结果1为成功，0为失败
		"message": "处理成功！",
		"data":    gin.H{},
	})
}

//照相
func CameraMonitor(dat map[string]interface{}, c *Client, sign string) {
	deviceID := webgo.GetResult(dat["deviceID"])
	data := dat["data"].(map[string]interface{})
	imageUrl := webgo.GetResult(data["imageUrl"])
	redis := db.CfRedis[0]
	android_screen_id := redis.HGet("device:id", deviceID).Val()
	var androidScreenPhoto model.Android_screen_photo
	asId, _ := strconv.Atoi(android_screen_id)
	androidScreenPhoto.Android_screen_id = asId
	androidScreenPhoto.Img = imageUrl
	androidScreenPhoto.Type = 2
	advertService.InsertScreenPhoto(androidScreenPhoto)
	c.Write(gin.H{
		"sign":    sign,
		"status":  1, //处理结果1为成功，0为失败
		"message": "处理成功！",
		"data":    gin.H{},
	})
}

//设备信息
func DeviceInfo(dat map[string]interface{}, c *Client, sign string) {
	deviceID := webgo.GetResult(dat["deviceID"])
	//data := dat["data"].(map[string]interface{})
	//appVersion := webgo.GetResult(data["appVersion"])
	//fmt.Print(appVersion)
	redis := db.CfRedis[0]
	id := redis.HGet("device:id", deviceID).Val() //:id为获取androidScreenId  :number 获取AndroidNumber  这里的deviceID为Androidnumber
	DeviceInfoResult := advertService.FindAndroidScreenAndDevice(id)
	var DeviceModulars []map[string]interface{}
	for _, val := range DeviceInfoResult.CommLgz {
		var temp map[string]interface{}
		temp = make(map[string]interface{})
		temp["id"] = val.ID
		temp["lgz_id"] = val.LgzID
		temp["install_site"] = val.DeviceModular.InstallSite
		temp["modular_type"] = val.DeviceModular.ModularType
		temp["status"] = val.Status
		temp["comm_number"] = val.LgzID
		temp["coin"] = val.DeviceModular.Coin
		DeviceModulars = append(DeviceModulars, temp)
	}
	for  _,val := range DeviceInfoResult.WeiMaQi{
		var temp map[string]interface{}
		temp = make(map[string]interface{})
		temp["id"] = val.ID
		temp["install_site"] = val.DeviceModular.InstallSite
		temp["modular_type"] = val.DeviceModular.ModularType
		temp["status"] = val.Status
		temp["comm_number"] = val.WeimaqiId
		temp["coin"] = val.DeviceModular.Coin
		DeviceModulars = append(DeviceModulars, temp)
	}
	if DeviceInfoResult.AndroidScreen.ID != 0 {
		c.Write(gin.H{
			"sign":       sign,
			"socketType": "deviceInfo",
			"deviceID":   deviceID,
			"status":     1, //处理结果1为成功，0为失败
			"message":    nil,
			"data": gin.H{
				"id":              DeviceInfoResult.Device.ID,
				"hardware_id":     DeviceInfoResult.Device.HardwareID,
				"model_id":        DeviceInfoResult.Device.ModelID,
				"is_bind":         DeviceInfoResult.Device.IsBind,
				"release_date":    DeviceInfoResult.Device.ReleaseDate,
				"store_id":        DeviceInfoResult.Device.StoreID,
				"status":          DeviceInfoResult.Device.Status,
				"device_number":   DeviceInfoResult.Device.DeviceNumber,
				"update_time":     DeviceInfoResult.Device.UpdateTime,
				"create_time":     DeviceInfoResult.Device.CreateTime,
				"qrcode":          DeviceInfoResult.QrCode,
				"device_modulars": DeviceModulars,
			},
		})
	} else {
		c.Write(gin.H{
			"status":     0, //处理结果1为成功，0为失败
			"sign":       sign,
			"socketType": "deviceInfo",
			"deviceID":   deviceID,
			"message":    "无数据！",
		})
	}
}

//新接口 暂未开发
func AndroidGetVisitorTotal(dat map[string]interface{}, c *Client, sign string) {
	webgo.GetResult(dat["socketType"])
	webgo.GetResult(dat["deviceID"])
	webgo.GetResult(dat["status"])
	data := dat["data"].(map[string]interface{})
	webgo.GetResult(data["todayTotal"])
	webgo.GetResult(data["todayNewTotal"])
	webgo.GetResult(data["todayStayTime"])
	webgo.GetResult(data["historyTotal"])
	webgo.GetResult(data["historyNewTotal"])
	webgo.GetResult(data["historyStayTime"])
}

//设置androidNumber和androidScreenID的对应关系
func SetIdToNumber() {
	redis := db.CfRedis[0]
	androidSreen := advertService.FindAllAS()
	for _, val := range androidSreen {
		redis.HSet("device:id", val.AndroidNumber, val.ID)                   //androidNum域存AndroidScreenID
		redis.HSet("device:number", strconv.Itoa(val.ID), val.AndroidNumber) //AndroidScreenID域存AndroidNumber
		//fmt.Println(val.AndroidNumber)
		redis.Del(val.AndroidNumber)
		//fmt.Println(redis.Exists(val.AndroidNumber).Val())
	}
}

//设置为离线状态
func SetOnlinStatus() {
	db.SqlDB.Exec("update android_screen set status=2")
}

//设置为在线状态
func online(deviceId string) {
	redis := db.CfRedis[0]
	//flag := redis.Get(deviceId).Val() //这边的deviceId为AndroidNumber
	if redis.Exists(deviceId).Val() ==0 {
		advertService.AndroidScreenStatus("1", deviceId) //设置为在线
	}
}

//	AdPushRedis 设置
func setRedisKey(key string, value string) {
	redis := db.CfRedis[0]
	maxAge := 8000
	duration_Milliseconds := time.Duration(maxAge) * time.Millisecond //过期时间8S
	err := redis.Set(key, value, duration_Milliseconds).Err()
	if err != nil {
		fmt.Println(err)
	}
}

//捕获异常抛出
func tryCathSend(c *Client) {
			c.Write(gin.H{
				"statue":0,
				"data": "请检查JSON数据",
			})
}

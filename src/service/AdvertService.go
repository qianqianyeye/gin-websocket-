package webservice

import (
	"AdPushServer_Go/src/db"
	"AdPushServer_Go/src/model"
	"fmt"
	"strings"
)

type AdverService struct {
}

func (ctrl *AdverService) GetAdById(id string) model.Advert {
	var s []model.Advert
	row := db.SqlDB.Table("advert").Select("*").Joins("left join advert_content on advert.id =advert_content.advert_id " +
		"left join advert_download on advert.id=advert_download.advert_id where advert.id=" + id).Scan(&s)
	fmt.Println(row.RowsAffected)
	//db.SqlDB.Joins("left join advert_content on advert.id =advert_content.advert_id").Joins("left join advert_download on advert.id=advert_download.advert_id where advert.id=?",id).Find(&s)
	var advert model.Advert
	var advertdownload []model.AdvertDownload
	var advertcontent []model.AdvertContent
	if len(s) > 0 {
		var temp []int
		db.SqlDB.Table("advert").Where("id=?", s[0].ID).Scan(&advert)
		db.SqlDB.Table("advert_download").Where("advert_id=?", s[0].ID).Scan(&advertdownload)
		db.SqlDB.Table("advert_content").Where("advert_id=?", s[0].ID).Scan(&advertcontent)
		if len(advertdownload) > 0 {
			for _, val := range advertdownload {
				temp = append(temp, val.AndroidScreenID)
			}
		}
		var androidscreen []model.AndroidScreen
		db.SqlDB.Table("android_screen").Where("id in (?)", temp).Scan(&androidscreen)
		advert.AdvertDownload = advertdownload
		advert.AdvertContent = advertcontent
		advert.AndroidScreen = androidscreen
	}
	return advert
}

func (ctrl *AdverService) Validata(AdvertDownload model.AdvertDownload) bool {
	androidScreenId := AdvertDownload.AndroidScreenID
	var deviceMoular model.DeviceModular
	db.SqlDB.Table("device_modular").Where("modular_id=?", androidScreenId).Scan(&deviceMoular)
	if deviceMoular.ID != 0 {
		var device model.Device
		db.SqlDB.Table("device").Where("id = ? ", deviceMoular.DeviceID).Scan(&device)
		if device.ID != 0 {
			var merchantStore model.MerchantStore
			db.SqlDB.Table("merchant_store").Where("store_id=?", device.StoreID).Scan(&merchantStore)
			if merchantStore.ID != 0 {
				if merchantStore.MerchantID == AdvertDownload.MerchantID {
					return true
				} else {
					db.SqlDB.Table("advert_download").Where("id=?", AdvertDownload.ID).Update("download_progress", "-1")
					return false
				}
			}
		}
	}
	return false
}

func (ctrl *AdverService) GetSocketIdKey(arr []int) []string {
	var androidscreen []model.AndroidScreen
	db.SqlDB.Table("android_screen").Where("id in (?)", arr).Scan(&androidscreen)
	var SocketKey = []string{}
	for _, val := range androidscreen {
		SocketKey = append(SocketKey, val.AndroidNumber)
	}
	return SocketKey
}

//获取所有android_screen
func (ctrl *AdverService) FindAllAS() []model.AndroidScreen {
	var androidScreens []model.AndroidScreen
	db.SqlDB.Table("android_screen").Scan(&androidScreens)
	return androidScreens
}
func (ctrl *AdverService) FindAllASByNum(androidNum string) []model.AndroidScreen {
	var androidScreens []model.AndroidScreen
	db.SqlDB.Table("android_screen").Where("android_number=?", androidNum).Scan(&androidScreens)
	return androidScreens
}

//根据屏幕id找到最新的一条广告
func (ctrl *AdverService) FindAdnroidScreenOne(androidScreenId string) model.AdvertDownload {
	var adverDownLoad model.AdvertDownload
	db.SqlDB.Table("advert_download").Where("android_screen_id=? and download_progress != ?", androidScreenId, "-1").Order("create_time desc limit 1").Scan(&adverDownLoad)
	return adverDownLoad
}

func (ctrl *AdverService) FindAndroidScreenAndDevice(androidScreenId string) model.DeviceInfoResult {

	var AndroidScreen model.AndroidScreen
	var DeviceModular model.DeviceModular
	var Device model.Device
	var Commlgz []model.CommLgz
	var WeiMaQi []model.CommWeimaqi
	var DeviceInfoResult model.DeviceInfoResult
	var zDeviceModular []model.DeviceModular
	var lgzId []int
	var wmqId []int
	/*
		根据Id获取Android_screen的数据，如果存在，则去查找device_modular关系表，获取到对应的设备ID,
		级联查询对应设备的乐关注信息
	*/
	db.SqlDB.Table("android_screen").Where("id=?", androidScreenId).Scan(&AndroidScreen)
	if AndroidScreen.ID != 0 {
		//modular_type :1：唯码器，2：安卓屏幕，3：乐关注 4:乐摇摇
		db.SqlDB.Table("device_modular").Where("modular_type=? and modular_id=?", "2", AndroidScreen.ID).Scan(&DeviceModular)
		if DeviceModular.ID != 0 {
				db.SqlDB.Table("device").Where("id=?",DeviceModular.DeviceID).Scan(&Device)
				db.SqlDB.Table("device_modular").Where("device_id=?",Device.ID).Scan(&zDeviceModular)
				for _,val :=range zDeviceModular {
					if val.ModularType==3 {
						lgzId=append(lgzId, val.ModularID)
					}
					if  val.ModularType==1{
						wmqId=append(wmqId, val.ModularID)
					}
				}
				db.SqlDB.Table("comm_lgz").Select("*").Where("id in (?)",lgzId).Scan(&Commlgz)
				db.SqlDB.Table("comm_weimaqi").Select("*").Where("id in (?)",wmqId).Scan(&WeiMaQi)
				for i,val :=range Commlgz {
					for _,zval :=range zDeviceModular {
						if val.ID==zval.ModularID {
							Commlgz[i].DeviceModular=zval
						}
					}
				}
				for i,val :=range WeiMaQi {
					for _,zval :=range zDeviceModular {
						if val.ID==int64(zval.ModularID) {
							WeiMaQi[i].DeviceModular=zval
						}
					}
				}

		}
		//将查询到的信息封到结果集中，方便后面解析
		DeviceInfoResult.AndroidScreen = AndroidScreen
		DeviceInfoResult.DeviceModular = DeviceModular
		DeviceInfoResult.Device = Device
		DeviceInfoResult.CommLgz = Commlgz
		DeviceInfoResult.WeiMaQi=WeiMaQi
		DeviceInfoResult.QrCode = model.Claw + "?device_number=" + Device.DeviceNumber
		for _,val := range model.QrcAndroidNum {
			fmt.Println(val)
			if Device.DeviceNumber == val {
				DeviceInfoResult.QrCode = model.Claw + "miniapp?device_number=" + Device.DeviceNumber
			}
		}
	}
	return DeviceInfoResult
}

//插入图片到数据库（截屏1/照相2）
func (ctrl *AdverService) InsertScreenPhoto(androidScreenPhoto model.Android_screen_photo) {
	rows := db.SqlDB.Exec("insert into android_screen_photo (android_screen_id,img,type) values (?,?,?)", androidScreenPhoto.Android_screen_id, androidScreenPhoto.Img, androidScreenPhoto.Type)
	if rows.RowsAffected > 0 {
		fmt.Println("插入android_screen_photo成功")
	} else {
		fmt.Println("插入android_screen_photo失败")
	}
}

//更新广告进度
func (ctrl *AdverService) UpdateAdDownLoadProgress(parmMap map[string]string) {
	for k, v := range parmMap {
		//k: adverId+"_"+androidScreenId  v: adDownloadProgress
		arr := strings.Split(k, "_")
		advertId := arr[0]
		androidScreenId := arr[1]
		adDownloadProgress := v
		db.SqlDB.Table("advert_download").Where("android_screen_id=? and advert_id=?", androidScreenId, advertId).Update("download_progress", adDownloadProgress)
		//db.SqlDB.Exec("update advert_download set  download_progress= ? where android_screen_id = ? and advert_id=?",adDownloadProgress,androidScreenId,advertId)
	}
}

func (ctrl *AdverService) AndroidScreenStatus(status string, deviceid string) {
	db.SqlDB.Exec("update android_screen set status = ? where android_number = ?", status, deviceid)
}

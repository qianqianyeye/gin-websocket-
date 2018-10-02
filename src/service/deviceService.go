package webservice

import (
	"AdPushServer_Go/src/db"
	"AdPushServer_Go/src/model"
	"math/rand"
	"strconv"
	"time"
	"AdPushServer_Go/src/webgo"
)

type DeviceService struct {
}

func (ctrl *DeviceService) CreateDevice(androidId string) model.AndroidScreen {
	defer func() {
		if e := recover(); e != nil {
			webgo.Error("CreateDevice：%s",e)
		}
	}()
	var androidScreen model.AndroidScreen
	row := db.SqlDB.Table("android_screen").Where("android_model=?", androidId).Scan(&androidScreen) //查询androidId是否存在
	if row.RowsAffected > 0 {
		return androidScreen //存在
	} else {
		androidNum := getAndroidNum() //创建AndroidNumber
		androidScreen.AndroidNumber = androidNum
		androidScreen.AndroidModel = androidId
		//rows :=db.SqlDB.Table("android_screen").Create(androidScreen)
		t := time.Now()
		rows := db.SqlDB.Exec("insert into android_screen (android_model,android_number,status,create_time) values (?,?,?,?)", androidId, androidNum, "1", t)
		if rows.RowsAffected > 0 {
			return androidScreen
		}
	}
	return androidScreen
}

//创建AndroidNumber
func getAndroidNum() string {
	defer func() {
		if e := recover(); e != nil {
			webgo.Error("getAndroidNum：%s",e)
		}
	}()
	rand.Seed(time.Now().Unix())
	i := rand.Intn(9000) + 10000
	androidNum := "S" + strconv.Itoa(i) //S+  1万到10万之间的随机数
	var androidScreen model.AndroidScreen
	//row :=db.SqlDB.Table("android_screen").Where("android_number= ?",androidNum)
	db.SqlDB.Table("android_screen").Where("android_number= ?", androidNum).Scan(androidScreen) //判断表中是否有这个AndroidNumber
	//存在则重新创建
	if androidScreen.ID != 0 {
		return getAndroidNum()
	}
	return androidNum
}

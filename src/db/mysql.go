package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"

	"AdPushServer_Go/src/webgo"
	"time"
)

var SqlDB *gorm.DB

func InitDB(flag bool) {
	var err error
	//root:root@tcp(192.168.1.87:3306)/test?parseTime=true
	//SqlDB, err = gorm.Open("mysql", "SaasSH:mengdao+1ssh@tcp(119.23.219.245:8001)/SassDatabaseDev?charset=utf8")
	if flag {
		webgo.Debug("使用线上数据库%s","...")
		SqlDB, err = gorm.Open("mysql", "SaasSH:mengdao+1ssh@tcp(119.23.219.245:8001)/SassDatabase?charset=utf8&parseTime=True&loc=Local")
	}else{
		webgo.Debug("使用线下数据库%s","...")
		SqlDB, err = gorm.Open("mysql", "SaasSH:mengdao+1ssh@tcp(119.23.219.245:8001)/SassDatabaseDev?charset=utf8&parseTime=True&loc=Local")
	}
	SqlDB.LogMode(true)
	SqlDB.DB().SetMaxOpenConns(1000)
	SqlDB.DB().SetMaxIdleConns(500)
	SqlDB.DB().SetConnMaxLifetime(60*time.Second)
	SqlDB.DB().Ping()
	//SqlDB.
	if err != nil {
		log.Fatal(err.Error())
	}
}

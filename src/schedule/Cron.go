package schedule

import (
	"AdPushServer_Go/src/db"
	"AdPushServer_Go/src/service"
	"github.com/robfig/cron"
	"strconv"
	"strings"
)

var advertService webservice.AdverService

//广告进度定时器
func UpdateAdModelCronRun() {
	c := cron.New()
	//每隔30秒执行
	c.AddFunc("*/30 * * * * *", func() {
		//fmt.Println("更新广告进度中...")
		redis := db.CfRedis[0]
		lens := redis.LLen("updateAdModel").Val()                   //获取LLen数组的长度
		updateModel := redis.LRange("updateAdModel", 0, lens).Val() //取出所有
		var tempMap map[string]string
		tempMap = make(map[string]string)
		//遍历数组
		for _, val := range updateModel {
			//解析 adUpdateID+","+adDownloadProgress+","+android_screen_id
			arr := strings.Split(val, ",")
			adverId := arr[0]
			adDownloadProgress := arr[1]
			androidScreenId := arr[2]
			//判断tempMap中设备对应的广告是否存在，存在则比较进度，存最大的进度值
			//不存在则放入
			if _, ok := tempMap[adverId+"_"+androidScreenId]; ok {
				tempNow, _ := strconv.Atoi(adDownloadProgress)
				Old := tempMap[adverId+"_"+androidScreenId]
				tempOld, _ := strconv.Atoi(Old)
				if tempNow > tempOld {
					tempMap[adverId+"_"+androidScreenId] = adDownloadProgress
				}
			} else {
				tempMap[adverId+"_"+androidScreenId] = adDownloadProgress
			}
		}
		//更新
		if len(tempMap) > 0 {
			//经过比较后tempMap中存放的均为对应设备广告的最大下载进度
			advertService.UpdateAdDownLoadProgress(tempMap)
		}
		//清空redis中的updateAdModel数组
		redis.LTrim("updateAdModel", lens, -1)
		//fmt.Println("更新完毕！")
	})
	c.Start()
}

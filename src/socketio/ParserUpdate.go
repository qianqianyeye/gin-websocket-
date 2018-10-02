package socketio

import (
	"AdPushServer_Go/src/db"
	"AdPushServer_Go/src/model"
	"github.com/gin-gonic/gin"
)

//下发广告
func ParseAdJson(s *model.Advert, key []string, status string) {
	var adAreaData []map[string]interface{}
	if len(key) > 0 {
		//0自定义模板  其他为通用模板
		if s.Model == 0 {
			//自定义模板
			if len(s.AdvertContent) > 0 {
				for i := 0; i < len(s.AdvertContent)-1; i++ {
					maps := map[string]interface{}{}
					data := []map[string]interface{}{}
					f := false
					ad := s.AdvertContent[i]
					maps["areaWidth"] = ad.AreaWidth
					maps["areaHeight"] = ad.AreaHeight
					maps["x"] = ad.X
					maps["y"] = ad.Y
					for j := i + 1; j <= len(s.AdvertContent)-1; j++ {
						if ad.Flag {
							break
						}
						if ad.X == s.AdvertContent[j].X && ad.Y == s.AdvertContent[j].Y {
							tempUrl1 := map[string]interface{}{}
							tempUrl2 := map[string]interface{}{}
							tempUrl1["url"] = model.Imgurl + ad.URL
							data = append(data, tempUrl1)
							s.AdvertContent[i].Flag = true
							tempUrl2["url"] = model.Imgurl + s.AdvertContent[j].URL
							data = append(data, tempUrl2)
							s.AdvertContent[j].Flag = true
							f = true
						}
					}
					if f == true {
						maps["datas"] = data
						adAreaData = append(adAreaData, maps)
					}
				}
			}
			for _, val := range s.AdvertContent {
				//不同位置的URL
				if val.Flag == false {
					data := []map[string]interface{}{}
					tempUrl := map[string]interface{}{}
					maps := map[string]interface{}{}
					maps["areaWidth"] = val.AreaWidth
					maps["areaHeight"] = val.AreaHeight
					maps["x"] = val.X
					maps["y"] = val.Y
					tempUrl["url"] = model.Imgurl + val.URL
					data = append(data, tempUrl)
					maps["datas"] = data
					adAreaData = append(adAreaData, maps)
				}
			}

			for _, val := range key {
				redis := db.CfRedis[0]
				//获取
				clientId := redis.Get(val).Val() //androidnum
				//判断连接是否存在
				if _, ok := AndroidServer.Clients[clientId]; ok {
					result :=make(map[string]interface{})
					result["clientId"]=clientId
					result["data"]=gin.H{
						"status":     "1", //处理结果1为成功，0为失败
						"socketType": status,
						"deviceID":   val,
						"sign":       GetSign(val, status),
						"data": gin.H{
							"adModelType": s.Model,
							"scrollTime":  s.ScrollTime,
							"adUpdateID":  s.ID,
							"adAreaData":  adAreaData,
						},
					}
					AndroidServer.Msg(result)

				}
			}
		} else {
			for _, val := range s.AdvertContent {
				var url map[string]interface{} /*创建集合 */
				url = make(map[string]interface{})
				url["url"] = model.Imgurl + val.URL
				url["areaType"] = val.AreaType
				adAreaData = append(adAreaData, url)
			}

			for _, val := range key {
				redis := db.CfRedis[0]
				//获取
				clientId := redis.Get(val).Val() //androidNum
				//fmt.Println(redis.Get(val))
				//判断连接是否存在
				if _, ok := AndroidServer.Clients[clientId]; ok {
					result :=make(map[string]interface{})
					result["clientId"]=clientId
					result["data"]=gin.H{
						"socketType": status,
						"deviceID":   val,
						"status":     "1", //处理结果1为成功，0为失败
						"sign":       GetSign(val, status),
						"data": gin.H{
							"adModelType": s.Model,
							"scrollTime":  s.ScrollTime,
							"adUpdateID":  s.ID,
							"adAreaData":  adAreaData,
						},
					}
					AndroidServer.Msg(result)
				}
			}
		}
	}
}

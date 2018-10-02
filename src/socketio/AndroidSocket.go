package socketio

import (
	"AdPushServer_Go/src/db"
	"AdPushServer_Go/src/webgo"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func AndroidStartGame() {

}
func AndroidCtrlBtnPress(dat map[string]interface{}, c *Client, sign string){
	if webgo.GetResult(dat["status"]) != "0" {
		data := dat["data"].(map[string]interface{})
		lgzId := webgo.GetResult(data["grabBoardDeviceID"])
		status := webgo.GetResult(dat["status"])
		s,_:=strconv.Atoi(status)
		key := "claw_socket_" + lgzId
		redis := db.CfRedis[0]
		clientId := redis.Get(key).Val()
		result :=make(map[string]interface{})
		result["clientId"]=clientId
		result["data"]=gin.H{
			"socketType": "ctrlBtnPress",
			"status":   s,
			"data":    nil,
		}
		ClawServer.Msg(result)
		//if _, ok := ClawServer.Clients[clientId]; ok {
		//	ClawServer.Clients[clientId].Write(gin.H{
		//		"socketType": "ctrlBtnPress",
		//		"status":   s,
		//		"data":       nil,
		//	})
		//}
	}
}
func AndroidCatchBtnPress(dat map[string]interface{}, c *Client, sign string) {
	if webgo.GetResult(dat["status"]) != "0" {
		var data map[string]interface{}
		data = make(map[string]interface{})
		if dat["data"] != nil {
			data = dat["data"].(map[string]interface{})
		}
		lgzId := webgo.GetResult(data["grabBoardDeviceID"])
		status := webgo.GetResult(dat["status"])
		s,_:=strconv.Atoi(status)
		key := "claw_socket_" + lgzId
		redis := db.CfRedis[0]
		clientId := redis.Get(key).Val()

		result :=make(map[string]interface{})
		result["clientId"]=clientId
		result["data"]=gin.H{
			"socketType": "catchBtnPress",
			"status":     s,
			"data":       nil,
		}
		ClawServer.Msg(result)

		//if _, ok := ClawServer.Clients[clientId]; ok {
		//	ClawServer.Clients[clientId].Write(gin.H{
		//		"socketType": "catchBtnPress",
		//		"status":     s,
		//		"data":       nil,
		//	})
		//}
	}
}
func AndroidGameOver(dat map[string]interface{}, c *Client, sign string) {
	if webgo.GetResult(dat["status"]) != "0" {
		data := dat["data"].(map[string]interface{})
		lgzId := webgo.GetResult(data["grabBoardDeviceID"])
		status := webgo.GetResult(dat["status"])
		key := "claw_socket_" + lgzId
		redis := db.CfRedis[0]
		clientId := redis.Get(key).Val()
		gameStr := redis.Get("game_lgz_id_" + lgzId).Val()
		if gameStr!="" {
			var mdat map[string]interface{}
			if err := json.Unmarshal([]byte(gameStr), &mdat); err == nil {
				fmt.Println(mdat)
			} else {
				webgo.Error("AndroidGameOver json.Unmarshal failed:%s",err)
			}
			if _, ok := mdat["surplus_count"]; ok {
				temp := webgo.GetResult(mdat["surplus_count"])
				count, _ := strconv.Atoi(temp)
				count = count - 1
				mdat["surplus_count"] = count
				parm, err := json.Marshal(mdat)
				if err != nil {
					webgo.Error("json.Marshal failed:%s",err)
					//fmt.Println("json.Marshal failed:", err)
				}
				times, _ := strconv.Atoi(webgo.GetResult(mdat["time"]))
				hs := count*(times+5)*1000 + 5*1000
				duration_Milliseconds := time.Duration(hs) * time.Millisecond
				//	fmt.Println(duration_Milliseconds)
				redis.Set("game_lgz_id_"+lgzId, string(parm), duration_Milliseconds)
			}
			s,_:=strconv.Atoi(status)

			result :=make(map[string]interface{})
			result["clientId"]=clientId
			result["data"]=gin.H{
				"status":   s  ,
				"data":       mdat,
				"socketType": "gameOver",
			}
			ClawServer.Msg(result)

			//if _, ok := ClawServer.Clients[clientId]; ok {
			//	ClawServer.Clients[clientId].Write(gin.H{
			//		"status":   s  ,
			//		"data":       mdat,
			//		"socketType": "gameOver",
			//	})
			//}
		}
	}
}

func AndroidGetGrabTime(dat map[string]interface{}, c *Client, sign string) {
	if webgo.GetResult(dat["status"]) != "0" {
		data := dat["data"].(map[string]interface{})
		lgzId := webgo.GetResult(data["grabBoardDeviceID"])
		time := webgo.GetResult(data["time"])
		key := "claw_socket_" + lgzId
		redis := db.CfRedis[0]
		clientId := redis.Get(key).Val()

		result :=make(map[string]interface{})
		result["clientId"]=clientId
		result["data"]=gin.H{
			"socketType": "gameTime",
			"status":     1,
			"data": gin.H{
				"time": time,
			},
		}
		ClawServer.Msg(result)

		//if _, ok := ClawServer.Clients[clientId]; ok {
		//	ClawServer.Clients[clientId].Write(gin.H{
		//		"socketType": "gameTime",
		//		"status":     1,
		//		"data": gin.H{
		//			"time": time,
		//		},
		//	})
		//}
	}
}

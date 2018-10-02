package socketio

import (
	"AdPushServer_Go/src/db"
	"AdPushServer_Go/src/webgo"
	"github.com/gin-gonic/gin"
)

func ClawStartGame() {

}

func clawCtrlBtnPress(parmMap map[string]interface{}, dat map[string]interface{}) {
	androidNumber := webgo.GetResult(parmMap["androidNumber"])
	redis := db.CfRedis[0]
	clientId := redis.Get(androidNumber).Val()
	if clientId == "" {
		return
	}
	parmMap["grabBoardDeviceID"] = webgo.GetResult(parmMap["lgzId"])
	temp := dat["data"].(map[string]interface{})
	temp["grabBoardDeviceID"] = webgo.GetResult(parmMap["lgzId"])

	result :=make(map[string]interface{})
	result["clientId"]=clientId
	result["data"]=gin.H{
		"status":1,
		"socketType": "ctrlBtnPress",
		"deviceID":   androidNumber,
		"data":       temp,
		"sign":       GetSign(androidNumber, "ctrlBtnPress"),
	}
	AndroidServer.Msg(result)
	//if _, ok := AndroidServer.Clients[clientId]; ok {
	//	AndroidServer.Clients[clientId].Write(gin.H{
	//		"status":1,
	//		"socketType": "ctrlBtnPress",
	//		"deviceID":   androidNumber,
	//		"data":       temp,
	//		"sign":       GetSign(androidNumber, "ctrlBtnPress"),
	//	})
	//} else {
	//	return
	//}
}

func clawCatchBtnPress(parmMap map[string]interface{}, dat map[string]interface{}) {
	androidNumber := webgo.GetResult(parmMap["androidNumber"])
	redis := db.CfRedis[0]
	clientId := redis.Get(androidNumber).Val()
	if clientId == "" {
		return
	}
	//temp := dat["data"].(map[string]interface{})
	//temp["grabBoardDeviceID"]=webgo.GetResult(parmMap["lgzId"])

	result :=make(map[string]interface{})
	result["clientId"]=clientId
	result["data"]=gin.H{
		"status":     1,
		"socketType": "catchBtnPress",
		"deviceID":   androidNumber,
		"data": gin.H{
			"grabBoardDeviceID": parmMap["lgzId"],
		},
		"sign": GetSign(androidNumber, "catchBtnPress"),
	}
	AndroidServer.Msg(result)

	//if _, ok := AndroidServer.Clients[clientId]; ok {
	//	AndroidServer.Clients[clientId].Write(gin.H{
	//		"status":     1,
	//		"socketType": "catchBtnPress",
	//		"deviceID":   androidNumber,
	//		"data": gin.H{
	//			"grabBoardDeviceID": parmMap["lgzId"],
	//		},
	//		"sign": GetSign(androidNumber, "catchBtnPress"),
	//	})
	//} else {
	//	return
	//}
}

func clawGameTime(parmMap map[string]interface{}, dat map[string]interface{}) {
	androidNumber := webgo.GetResult(parmMap["androidNumber"])
	redis := db.CfRedis[0]
	clientId := redis.Get(androidNumber).Val()
	if clientId == "" {
		return
	}

	result :=make(map[string]interface{})
	result["clientId"]=clientId
	result["data"]=gin.H{
		"status":     1,
		"socketType": "getGrabTime",
		"deviceID":   androidNumber,
		"data": gin.H{
			"grabBoardDeviceID": parmMap["lgzId"],
		},
		"sign": GetSign(androidNumber, "getGrabTime"),
	}
	AndroidServer.Msg(result)

	//if _, ok := AndroidServer.Clients[clientId]; ok {
	//	AndroidServer.Clients[clientId].Write(gin.H{
	//		"status":     1,
	//		"socketType": "getGrabTime",
	//		"deviceID":   androidNumber,
	//		"data": gin.H{
	//			"grabBoardDeviceID": parmMap["lgzId"],
	//		},
	//		"sign": GetSign(androidNumber, "getGrabTime"),
	//	})
	//} else {
	//	return
	//}
}

func getVisitorTotal(dat map[string]interface{}) {
	androidNumber := webgo.GetResult(dat["deviceID"])
	socketType := webgo.GetResult(dat["socketType"])
	redis := db.CfRedis[0]
	clientId := redis.Get(androidNumber).Val()
	result :=make(map[string]interface{})
	result["clientId"]=clientId
	result["data"]=gin.H{
		"socketType": socketType,
		"deviceID":   androidNumber,
		"sign":      GetSign(androidNumber, socketType),
		"status": 1,
	}
	AndroidServer.Msg(result)
	//if _, ok := AndroidServer.Clients[clientId]; ok {
	//	AndroidServer.Clients[clientId].Write(gin.H{
	//		"socketType": socketType,
	//		"deviceID":   androidNumber,
	//		"sign":      GetSign(androidNumber, socketType),
	//		"status": 1,
	//	})
	//}
}

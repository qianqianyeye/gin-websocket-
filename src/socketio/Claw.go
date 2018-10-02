package socketio

import (
	"AdPushServer_Go/src/db"
	"AdPushServer_Go/src/model"
	"AdPushServer_Go/src/webgo"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"strings"
	"time"
)

func (c *Client) ClawListen(ctx *gin.Context) {
	go c.ClawListenWrite(ctx)
	c.ClawListenRead(ctx)
}

func (c *Client) ClawConn() *websocket.Conn {
	return c.Ws
}

func (c *Client) ClawDone() {
	c.DoneCh <- true
}

func (c *Client) ClawWrite(msg map[string]interface{}) {
	defer webgo.TryCatch()
	select {
	case c.Msg <- msg:
	default:
		c.Server.Del(c)
		//c.Server.DelClaw(c)
		err := fmt.Errorf("claw %d is disconnected.", c.Id)
		c.Server.Err(err)
	}
}

func (c *Client) ClawListenWrite(ctx *gin.Context) {
	defer webgo.TryCatch()
	for {
		select {
		case msg := <-c.Msg:
			//log.Println("ClawSend:", msg)
			websocket.JSON.Send(c.Ws, msg)
		case <-c.DoneCh:
			c.Server.Del(c)
			//c.Server.DelClaw(c)
			c.DoneCh <- true
			return
		}
	}
}

func (c *Client) ClawListenRead(ctx *gin.Context) {
	defer webgo.TryCatch()
	for {
		select {
		case <-c.DoneCh:
			c.Server.Del(c)
			//c.Server.DelClaw(c)
			c.DoneCh <- true
			return
		default:
			var msg string
			err := websocket.Message.Receive(c.Ws, &msg)
			if err == io.EOF {
				c.DoneCh <- true
			} else if err != nil {
				c.DoneCh <- true
				c.Server.Err(err)
			} else {
				dat := PaserStringToMap(msg)
				socketType := webgo.GetResult(dat["socketType"])
				flag, resultMap := AuthMiddle(ctx, c)
				if flag {
					switch socketType {
					case "startGame":
						ClawStartGame()
					case "ctrlBtnPress":
						clawCtrlBtnPress(resultMap, dat)
					case "catchBtnPress":
						clawCatchBtnPress(resultMap, dat)
					case "gameTime":
						clawGameTime(resultMap, dat)
					//获取当前设备客流量数据
					case "getVisitorTotal":
						getVisitorTotal(dat)
					}
				}
			}
		}
	}
}

func AuthMiddle(ctx *gin.Context, client *Client) (bool, map[string]interface{}) {
	var resultMap map[string]interface{}
	resultMap = make(map[string]interface{})
	ws := client.Ws
	token := ctx.DefaultQuery("token", "")
	if token == "" {
		websocket.JSON.Send(ws, gin.H{
			"status": 0,
			"data":   "token数据为空！",
		})
		return false, resultMap
	}
	temp := strings.Split(token, " ")
	mySignKeyBytes := []byte(model.ClawKey)
	parseAuth, err := jwt.Parse(temp[1], func(*jwt.Token) (interface{}, error) {
		return mySignKeyBytes, nil
	})
	if err != nil {
		websocket.JSON.Send(ws, gin.H{
			"status": 0,
			"data":   "非法签名！",
		})
		client.Server.Del(client)
		//client.Server.DelClaw(client)
		client.DoneCh <- true
		log.Println("Now Claw 连接断开")
		return false, resultMap
	}
	claim := parseAuth.Claims.(jwt.MapClaims)
	var parmMap map[string]interface{}
	parmMap = make(map[string]interface{})
	for key, val := range claim {
		parmMap[key] = val
	}
	lgzId := ctx.Query("lgz_id")
	redis := db.CfRedis[1]
	data := parmMap["data"].(map[string]interface{})
	memberId := webgo.GetResult(data["id"])
	nonce := redis.Get(memberId).Val()
	nonceStr := webgo.GetResult(parmMap["nonceStr"])
	if nonceStr != nonce {
		websocket.JSON.Send(ws, gin.H{
			"status": 0,
			"data":   "noceStr错误！",
		})
		return false, resultMap
	}
	adRedis := db.CfRedis[0]
	gameStr := adRedis.Get("game_lgz_id_" + lgzId).Val()
	if gameStr == "" {
		websocket.JSON.Send(ws, gin.H{
			"status": 0,
			"data":   "获取数据失败！",
		})
		return false, resultMap
	}
	mapGameStr := PaserStringToMap(gameStr)
	if webgo.GetResult(mapGameStr["member_id"]) != memberId {
		websocket.JSON.Send(ws, gin.H{
			"status": 0,
			"data":   "当前状态不可游戏！",
		})
		return false, resultMap
	} else {
		resultMap["lgzId"] = lgzId
		resultMap["memberId"] = memberId
		resultMap["androidNumber"] = mapGameStr["android_number"]
		redis.Set("claw_socket_"+lgzId, client.Id, time.Hour)
		return true, resultMap
	}
}

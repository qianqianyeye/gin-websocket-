package socketio

import (
	"AdPushServer_Go/src/db"
	"AdPushServer_Go/src/webgo"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"io"
	"time"
)

type Client struct {
	Id     string
	Msg    chan interface{}
	Ws     *websocket.Conn
	Server *Server
	DoneCh chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	defer func() {
		if e := recover(); e != nil {
			webgo.Error("NewClient！ %s",e)
		}
	}()
	if ws == nil {
		panic("ws cannot be nil")
	}
	if server == nil {
		panic("server cannot be nil")
	}
	ID := uuid.Must(uuid.NewV4()).String()
	doneCh := make(chan bool)
	msg := make(chan interface{})
	return &Client{ID, msg, ws, server, doneCh}
}

func (c *Client) AndroidListen() {
	flag := make(chan string,1)
	go OverTime(flag,c)
	go c.AndroidListenWrite()
	c.AndroidListenRead(flag)
}

func (c *Client) Conn() *websocket.Conn {
	return c.Ws
}

func (c *Client) Done() {
	c.DoneCh <- true
}

//服务端写入数据到对应Client的chan中
func (c *Client) Write(msg interface{}) {
	select {
	case c.Msg <- msg:
	default:
		c.Server.Del(c)
		//c.Server.DelAndroid(c)
		fmt.Errorf("有错误发生了！s%", "断开连接")
		//c.Server.Err(err)
	}
}

//监听对应客户端的chan信道的写入 有写入则发送
func (c *Client) AndroidListenWrite() {
//	log.Println("Listening write to client")
	for {
		select {
		case msg := <-c.Msg:
			//log.Println("Send:", msg)
			websocket.JSON.Send(c.Ws, msg)
		case <-c.DoneCh:
			//断开删除对应的客户
			c.Server.Del(c)
			//c.Server.DelAndroid(c)
			c.DoneCh <- true
			//fmt.Print(c.DoneCh)
			return
		}
	}
}

//Android访问事件
func (c *Client) AndroidListenRead(f chan string) {
	for {
		select {
		case <-c.DoneCh:
			c.Server.Del(c)
			c.DoneCh <- true
			return
		default:
			var msg string
			err := websocket.Message.Receive(c.Ws, &msg) //接收发送方信息
			if err == io.EOF {
				c.DoneCh <- true
			} else if err != nil {
				c.DoneCh <- true
				webgo.Error("websocket出错了！ %s",err)
				c.Server.Err(err)
			} else {
				dat := PaserStringToMap(msg) //解析成对应的Map
				socketType := webgo.GetResult(dat["socketType"])
				flag, sign := ValiSign(dat, c) //验证签名
				if webgo.GetResult(dat["status"]) == "0" {
					flag = false
				}
				if flag {
					if socketType != "heartImpulse" {
						db.SqlDB.DB().Ping() //解决有时候会报连接错误   可能是闲置太久导致的问题
					}
					switch socketType {
					case "appVersion":
						AppVersion(dat,c,sign)
					case "createDeviceID":
						CreateDeviceID(dat, c, sign)
					case "heartImpulse":
						HeartImpulse(dat, c, sign)
						f<-webgo.GetResult(dat["deviceID"])
					case "updateAdModel":
						UpdateModel(dat, c, sign)
					case "checkUpdateID":
						CheckUpdateID(dat, c, sign)
					case "screenShot":
						ScreenShot(dat, c, sign)
					case "cameraMonitor":
						CameraMonitor(dat, c, sign)
					case "deviceInfo":
						DeviceInfo(dat, c, sign)
					case "startGame":
						AndroidStartGame()
					case "ctrlBtnPress":
						AndroidCtrlBtnPress(dat, c, sign)
					case "catchBtnPress":
						AndroidCatchBtnPress(dat, c, sign)
					case "getGrabTime":
						AndroidGetGrabTime(dat, c, sign)
					case "gameOver":
						AndroidGameOver(dat, c, sign)
					case "getVisitorTotal":
						AndroidGetVisitorTotal(dat, c, sign)
					}
				}
			}
		}
	}
}

func OverTime(f chan string,c *Client)  {
	var androidNum string
	DONE:
	for  {
		select {
		case <-time.After(8 * time.Second):
			delete(AndroidServer.Clients, c.Id)
			db.SqlDB.DB().Ping()                                                                          //解决有时候会报连接错误   可能是闲置太久导致的问题
			db.SqlDB.Exec("update android_screen set status = ? where android_number=?", "2", androidNum) //设置为下线
			//fmt.Println("断开")
			break DONE
		case a:=<-f:
			androidNum=a
		}
	}
}

func PaserStringToMap(msg string) map[string]interface{} {
	//字符串JSON转map
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(msg), &dat); err == nil {

	} else {
		webgo.Error("paserStringtoMap解析出错！")
	}
	return dat
}
func ValiSign(dat map[string]interface{}, c *Client) (bool, string) {
	defer func() {
		if e := recover(); e != nil {
			webgo.Error("ValiSign:%s",e)
			tryCathSend(c)
		}
	}()
	socketType := webgo.GetResult(dat["socketType"])
	deviceId :=webgo.GetResult(dat["deviceID"])
	csign := webgo.GetResult(dat["sign"])
	var sign string
	if socketType == "createDeviceID" {
		sign = GetSign("", socketType)
	} else {
		deviceId := webgo.GetResult(dat["deviceID"])
		sign = GetSign(deviceId, socketType)
	}
	if csign != sign {
		if socketType == "createDeviceID" {
			c.Write(gin.H{
				"socketType": socketType,
				"status":     0,       //处理结果1为成功，0为失败
				"message":    "签名错误！", //失败时返回错误信息
				"sign":GetSign("",socketType),
			})
		}else {
			c.Write(gin.H{
				"socketType": socketType,
				"status":     0,       //处理结果1为成功，0为失败
				"message":    "签名错误！", //失败时返回错误信息
				"sign":GetSign( webgo.GetResult(dat["deviceID"]),socketType),
			})
		}
		return false, ""
	} else if socketType == "createDeviceID"{
		return true, sign
	}else if socketType=="" || deviceId ==""{
			c.Write(gin.H{
				"socketType": socketType,
				"status":     0,       //处理结果1为成功，0为失败
				"message":    "JSON数据错误！", //失败时返回错误信息
				"sign":GetSign( webgo.GetResult(dat["deviceID"]),socketType),
			})
		return false, ""
	}else {
		return true, sign
	}
}

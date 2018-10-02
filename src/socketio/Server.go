package socketio

import (
	"AdPushServer_Go/src/db"
	"AdPushServer_Go/src/model"
	"AdPushServer_Go/src/webgo"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"strings"
	"time"
)
type Server struct {
	Pattern string
	Clients map[string]*Client
	doneCh    chan bool
	curClient   chan *map[string]interface{}
	addCh     chan *Client
	delCh  chan *Client
	errCh     chan error
	send chan interface{}
	msg chan map[string]interface{}
	sendClientId chan string
}

func NewServer(pattern string) *Server {
	clients := make(map[string]*Client)
	errCh := make(chan error)
	doneCh := make(chan bool)
	curClient := make(chan *map[string]interface{})
	addCh := make(chan *Client)
	delCh  :=make(chan *Client)
	send := make(chan interface{})
	msg := make(chan map[string]interface{},1000)
	sendClientId := make(chan string,1000)
	return &Server{
		pattern,
		clients,
		doneCh,
		curClient,
		addCh,
		delCh,
		errCh,
		send,
		msg,
		sendClientId,
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server)Msg(result map[string]interface{}) {
	s.msg<-result
}

func (s *Server)ClawServerListen()  {
	for {
		select {
		case c := <-s.addCh:
			s.Clients[c.Id] = c //在Server中存入连接对象的信息  Map
			webgo.Debug("有新客户端端加入，当前连接数为 %d",len(s.Clients))

		case c := <-s.delCh:
			s.DelClaw(c)
			webgo.Debug("有客户端断开连接.. 当前连接数 %d",len(s.Clients))

		case msg := <-s.msg:
			clientId:=msg["clientId"].(string)
			if _, ok :=s.Clients[clientId]; ok {
				s.Clients[clientId].Write(msg["data"])
			}

		case err := <-s.errCh:
			webgo.Error("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}

func (s *Server)AndroidServerListen()  {
	for {
		select {
		case c := <-s.addCh:
			s.Clients[c.Id] = c //在Server中存入连接对象的信息  Map
			webgo.Debug("有新安卓端加入，当前连接数为 %d",len(s.Clients))

		case c := <-s.delCh:
			s.DelAndroid(c)
			webgo.Debug("有安卓端断开连接.. 当前连接数 %d",len(s.Clients))

		case msg := <-s.msg:
			clientId:=msg["clientId"].(string)
			if _, ok :=s.Clients[clientId]; ok {
				s.Clients[clientId].Write(msg["data"])
			}

		case err := <-s.errCh:
			webgo.Error("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}

func (s *Server) DelAndroid(c *Client) {
	//c.Ws.Close()
	delete(s.Clients, c.Id)
	redis := db.CfRedis[0]
	androidNum := redis.Get(c.Id).Val()
	redis.Del(androidNum)
	redis.Del(c.Id)
	db.SqlDB.DB().Ping()                                                                          //解决有时候会报连接错误   可能是闲置太久导致的问题
	db.SqlDB.Exec("update android_screen set status = ? where android_number=?", "2", androidNum) //设置为下线
}

//Claw端断开连接
func (s *Server) DelClaw(c *Client) {
	delete(s.Clients, c.Id)
}


func (s *Server) AndroidListen(ctx *gin.Context, ws *websocket.Conn) {
	client := NewClient(ws, s) //创建一个连接的客户端
	//fmt.Println(len(s.Clients))
	//fmt.Println("add Client")
	s.Add(client)
	client.AndroidListen() //Android端监听读写事件
}

func (s *Server) ClawListen(ctx *gin.Context, ws *websocket.Conn) {
	client := NewClient(ws, s)
	//fmt.Println("add Claw")
	s.Add(client)
	//log.Println("Now Claw", len(s.Clients), " connected.")
	webgo.Debug("有新客户端加入，当前连接数为 %d",len(s.Clients))
	flag := ClawValidate(ctx, ws, client) //连接验证
	//flag:=true
	if flag {
		client.ClawListen(ctx)
	}
}

//Claw端连接验证
func ClawValidate(ctx *gin.Context, ws *websocket.Conn, client *Client) bool {
	defer func() {
		if e := recover(); e != nil {
			websocket.JSON.Send(ws, gin.H{
				"data": "出错了！连接断开",
				"socketType":"gameStatus",
				"status":1,
			})
			webgo.Error("客户端连接验证出错了！",e)
			client.Server.Del(client)
			client.DoneCh <- true
		}
	}()
	tokens := ctx.DefaultQuery("token", "") //获取token
	adRedis := db.CfRedis[0]                //adPush
	redis := db.CfRedis[1]                  //clawGlad
	if tokens == "" {
		websocket.JSON.Send(ws, gin.H{
			"status": 0,
			"socketType":"gameStatus",
			"data":   "token数据为空！关闭连接..",
		})
		webgo.Error("客户端token数据为空！即将关闭连接..")
		client.Server.Del(client)
		client.DoneCh <- true
		return false
	}
	mySignKeyBytes := []byte(model.ClawKey) //token的Key值
	temp := strings.Split(tokens, " ")
	//验证token
	parseAuth, err := jwt.Parse(temp[1], func(*jwt.Token) (interface{}, error) {
		return mySignKeyBytes, nil
	})
	if err != nil {
		websocket.JSON.Send(ws, gin.H{
			"status": 0,
			"socketType":"gameStatus",
			"data":   "非法签名！",
		})
		webgo.Error("客户端非法签名！即将关闭连接..%s",err)
		client.Server.Del(client)
		client.DoneCh <- true
		//log.Println("Now Claw 连接断开")
		return false
		//fmt.Println("parase with claims failed.", err)
	}
	//将token中的内容存入parmMap
	claim := parseAuth.Claims.(jwt.MapClaims)
	var parmMap map[string]interface{}
	parmMap = make(map[string]interface{})
	for key, val := range claim {
		parmMap[key] = val
	}
	lgzId := ctx.Query("lgz_id")
	data := parmMap["data"].(map[string]interface{})
	memberId := webgo.GetResult(data["id"])
	nonce := redis.Get(memberId).Val()
	nonceStr := webgo.GetResult(parmMap["nonceStr"])
	if nonceStr != nonce {
		websocket.JSON.Send(ws, gin.H{
			"status": 0,
			"socketType":"gameStatus",
			"data":   "noceStr错误！",
		})
		webgo.Error("客户端noceStr错误！即将关闭连接..")
		//log.Println("Now Claw 连接断开")
		client.Server.Del(client)
		client.DoneCh <- true
		return false
	}
	gameStr := adRedis.Get("game_lgz_id_" + lgzId).Val()
	if gameStr == "" {
		websocket.JSON.Send(ws, gin.H{
			"status": 2,
			"socketType":"gameStatus",
			"msg":  "当前可以游戏！",
			"data": nil,
		})
		webgo.Debug("成功！当前可以游戏！")
		adRedis.Set("claw_socket_"+lgzId, client.Id, time.Hour)
		return true
	}
	mapGameStr := PaserStringToMap(gameStr)
	if _, ok := mapGameStr["member_id"]; ok {
		if webgo.GetResult(mapGameStr["member_id"]) != memberId {
			websocket.JSON.Send(ws, gin.H{
				"status": 0,
				"socketType":"gameStatus",
				"data":   "当前状态不可游戏！",
			})
			webgo.Error("当前状态不可游戏！，即将关闭连接..")
			client.Server.Del(client)
			client.DoneCh <- true
			return false
		} else {
			webgo.Debug("成功！初始化！")
			websocket.JSON.Send(ws, gin.H{
				"status": 1,
				"socketType":"gameStatus",
				"data":   mapGameStr,
			})
		}
	}
	adRedis.Set("claw_socket_"+lgzId, client.Id, 0)
	return true
}
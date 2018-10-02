package main

import (
	"AdPushServer_Go/src/controller"
	"AdPushServer_Go/src/db"
	"AdPushServer_Go/src/schedule"
	"AdPushServer_Go/src/service"
	"AdPushServer_Go/src/socketio"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"golang.org/x/net/websocket"
	"time"
	"flag"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"net/http"
	"AdPushServer_Go/src/webgo"
)

func registerRouter(router *gin.Engine) {
	new(controller.AdvertController).Router(router)
	new(controller.TestController).Router(router)
}

var advertService webservice.AdverService


func main() {
	//webgo.Configuration("src/config/log4g.xml")
	l4g.LoadConfiguration("config/log4g.xml") //使用加载配置文件,类似与java的log4j.propertites
	defer l4g.Close()               //注:如果不是一直运行的程序,请加上这句话,否则主线程结束后,也不会输出和log到日志文件
	dataBase := flag.Bool("MySql",false,"true :线上，false: 线下 默认:false")
	flag.Parse()
	//fmt.Println(*dataBase)
	//*dataBase=true
	db.InitDB(*dataBase) //初始化数据库
	db.InitRedis(*dataBase)                //初始化Redis
	defer db.SqlDB.Close()
	go schedule.UpdateAdModelCronRun() //启动定时器
	socketio.SetIdToNumber()           //DeviceID和AndroidNumber对应关系
	socketio.SetOnlinStatus()          //将所有连接设置为不在线 2

	router := gin.Default()
	//网页跨域问题
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	androidServer := socketio.InitAndroidWebSocket()
	clawServer := socketio.InitClawWebSocket()

	//监听安卓websocket连接
	router.GET("/android", func(c *gin.Context) {
	//	fmt.Println(c.ClientIP())
		c.Request.Header.Add("Origin", "http://localhost:8010")
		handler := websocket.Handler(func(conn *websocket.Conn) {
			//conn.Request().Header.Set("Access-Control-Allow-Origin", "*")
			defer webgo.TryCatch()
			androidServer.AndroidListen(c, conn)
		})
		handler.ServeHTTP(c.Writer, c.Request)
	})

	//监听claw websocket连接
	router.GET("/claw", func(c *gin.Context) {
		//fmt.Println("coming1")
		c.Request.Header.Add("Origin", "http://localhost:8010")
		//fmt.Println("coming1")
		handler := websocket.Handler(func(conn *websocket.Conn) {
			//fmt.Println("coming2")
			defer webgo.TryCatch()
			clawServer.ClawListen(c, conn)
		})
		handler.ServeHTTP(c.Writer, c.Request)
	})

	router.LoadHTMLGlob("asset/ts/*")
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", "")
	})
	router.GET("/test", func(context *gin.Context) {
		fmt.Println("run ...")
	})
	//注册Http请求
	registerRouter(router)
	router.Run(":8010")
}

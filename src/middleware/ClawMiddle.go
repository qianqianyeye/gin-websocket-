package middleware

import (
	"AdPushServer_Go/src/db"
	"AdPushServer_Go/src/socketio"
	"AdPushServer_Go/src/webgo"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

//验证
func ClawMiddle(c *gin.Context) {
	defer webgo.TryCatchWeb(c)
	buf, _ := c.GetRawData()
	var str string = string(buf[0:len(buf)])
	var rmap map[string]interface{}
	if err := json.Unmarshal([]byte(str), &rmap); err == nil {
	} else {
		fmt.Println(err)
	}
	webgo.Debug("接收到参数：%+v\n",rmap)
	fmt.Println("接收到参数：\n",rmap)
	var keys []string
	for k := range rmap {
		if k != "sign" && k != "timeStamp" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	var result string
	for _, val := range keys {
		tempStr := webgo.GetResult(rmap[val])
		c.Set(val, rmap[val])
		result = result + "&" + val + "=" + tempStr
	}
	result = string([]rune(result)[1:]) + webgo.GetResult(rmap["timeStamp"])
	webgo.Debug("加密前字符串：%+v\n",result)
	fmt.Println("加密前字符串：\n",result)
	serverSign := socketio.GetMd5(result)
	webgo.Debug("加密后字符串：%+v\n",serverSign)
	fmt.Println("加密后字符串：\n",serverSign)
	webgo.Debug("要对比的加密字符串：%+v\n",webgo.GetResult(rmap["sign"]))
	fmt.Println("要对比的加密字符串：\n",webgo.GetResult(rmap["sign"]))

	if serverSign != webgo.GetResult(rmap["sign"]) {
		c.Abort()
		c.JSON(http.StatusOK, gin.H{
			"data": "签名错误!",
			"status":1,
		})
		return
	}
	db.SqlDB.DB().Ping() //解决有时候会报连接错误   可能是闲置太久导致的问题
	c.Next()
}

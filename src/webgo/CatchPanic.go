package webgo

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

func TryCatch()  {
	if e := recover(); e != nil {
		if e := recover(); e != nil {
			var err error
			switch x := e.(type) {
			case error:
				err = x
			case string:
				err = errors.New(x)
			default:
				err = errors.New("UnKnow panic")
			}
			error := errors.Wrap(err, "")
			Error("%+v\n", error)
		}
	}
}

func TryCatchWeb(c *gin.Context)  {
	if e := recover(); e != nil {
		var err error
		switch x:=e.(type) {
		case error:
			err=x
		case string:
			err=errors.New(x)
		default:
			err=errors.New("UnKnow panic")
		}
		error :=errors.Wrap(err,"")
		Error("%+v\n",error)
		c.JSON(http.StatusOK, gin.H{"status": 1, "data": nil, "msg": "有错误发生了！请联系管理员！"})
	}
}
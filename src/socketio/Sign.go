package socketio

import (
	"crypto/md5"
	"encoding/hex"
)

const key string = "MengDaoAndroid"

func GetSign(deviceId string, socketType string) string {
	var result string
	if deviceId == "" {
		resulttemp := socketType + key
		result = GetMd5(resulttemp)
	} else {
		resulttemp := deviceId + socketType + key
		result = GetMd5(resulttemp)
	}
	return result
}

func GetMd5(str string) string {
	sign := md5.New()
	result := []byte(str)
	sign.Write(result)
	cipherStr := sign.Sum(nil)
	strResult := hex.EncodeToString(cipherStr)
	return strResult
}

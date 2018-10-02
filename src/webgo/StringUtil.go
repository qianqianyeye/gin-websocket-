package webgo

import (
	"strconv"
)

//转换为字符串类型
func GetResult(i interface{}) string {
	switch i.(type) {
	case string:
		return i.(string)
	case int:
		return strconv.Itoa(i.(int))
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(i.(float64), 'f', -1, 32)
	case []interface{}:
		temp := i.([]interface{})
		var str string
		for val := range temp {
			//fmt.Println(temp[val])
			s := GetResult(temp[val])
			str = str + "," + s
		}
		str = string([]rune(str)[1:])
		str =  str
		return str
	case interface{}:
		return i.(string)
	default:
		return ""
	}
}

func GetArr(i interface{}) []string {
	temp := i.([]interface{})
	var str []string
	for val := range temp {
		s := GetResult(temp[val])
		str = append(str, s)
	}
	return str
}

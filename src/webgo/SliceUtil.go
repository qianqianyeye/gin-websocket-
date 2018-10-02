package webgo

import "reflect"

func StringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	if (a == nil) != (b == nil) {
		return false
	}
	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// SliceEquals 用以比较两个Slice(基础数据类型,如[]int)内含值是否相等
func SliceEquals(a, b interface{}) bool {
	// a,有任意一个不是slice返回false
	_a := reflect.ValueOf(a)
	if _a.Kind() != reflect.Slice {
		panic("param a must be a slice")
	}
	_b := reflect.ValueOf(b)
	if _b.Kind() != reflect.Slice {
		panic("param a must be a slice")
	}
	// 长度不等则两个slice不同
	if _a.Len() != _b.Len() {
		return false
	}
	// 依次比较每个值
	for i := 0; i < _a.Len(); i++ {
		if _a.Index(i).Interface() != _b.Index(i).Interface() {
			return false
		}
	}
	return true
}
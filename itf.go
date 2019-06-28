package utils

import (
	"reflect"
	"time"
)

func IdxOfItfs(target interface{}, strs ...interface{}) int {
	for idx, str := range strs {
		if target == str {
			return idx
		}
	}
	return -1
}

// TODO 添加类型
// 检查值是否为默认空值
func IsBlank(v interface{}) bool {
	var blank interface{}
	switch v.(type) {
	case bool:
		blank = false

	case int:
		blank = 0

	case uint:
		blank = uint(0)

	case uint8:
		blank = uint8(0)

	case uint16:
		blank = uint16(0)

	case uint32:
		blank = uint32(0)

	case uint64:
		blank = uint64(0)

	case int64:
		blank = int64(0)

	case float32:
		blank = float32(0)

	case float64:
		blank = float64(0)

	case string:
		blank = ""

	case time.Time:
		blank = time.Time{}

	case *time.Time:
		blank = new(time.Time)

	default:
		// for interface
		if v == nil {
			return true
		}

		// for array
		if reflect.ValueOf(v).Len() == 0 {
			return true
		}

		return false
		//fmt.Println("the type %s can not support", v)
	}

	return v == blank
}

func Clone(i interface{}) interface{} {
	// Wrap argument to reflect.Value, dereference it and return back as interface{}
	//new := reflect.Indirect(reflect.ValueOf(i)).Addr().Interface()

	//return new
	/*	val := reflect.ValueOf(i)
		if val.CanAddr() {
			return val.Interface()
		}

		return val.Addr().Interface()
	*/

	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		// Pointer:
		return reflect.New(reflect.ValueOf(i).Elem().Type()).Interface()
	} else {
		// Not pointer:
		return reflect.New(reflect.TypeOf(i)).Elem().Interface()
	}
}

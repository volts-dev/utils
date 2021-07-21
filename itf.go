package utils

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func Itf2Bool(val interface{}) (res bool) {
	if val == nil {
		return
	}

	if value, ok := val.(bool); ok {
		return value
	} else {
		t := reflect.TypeOf(val)
		vv := reflect.Indirect(reflect.ValueOf(val))
		switch t.Kind() {
		case reflect.Bool:
			return vv.Bool()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return vv.Int() != 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return vv.Uint() != 0
		case reflect.String:
			if b, err := strconv.ParseBool(vv.String()); err != nil {
				fmt.Printf("Unsupported type Itf2Bool(%v) error : %s", vv.Type().Name(), err.Error())
			} else {
				return b
			}
		default:
			fmt.Printf("Unsupported type Itf2Bool(%v)", vv.Type().Name())
		}
	}

	return
}

func Itf2Int(val interface{}) (res int) {
	if val == nil {
		return
	}

	if value, ok := val.(int); ok {
		return value
	} else {
		t := reflect.TypeOf(val)
		vv := reflect.Indirect(reflect.ValueOf(val))
		switch t.Kind() {
		//checked
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int(vv.Int())

			//checked
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int(vv.Uint())

		//checked
		case reflect.Float32, reflect.Float64:
			return int(vv.Float())

			//checked
		case reflect.String:
			return StrToInt(vv.String())
			/*
				if i, err := strconv.ParseInt(vv.String(), 10, 0); err != nil {
					fmt.Printf("Unsupported type %v", vv.Type().Name())
				} else {
					return i
				}
			*/
		case reflect.Array, reflect.Slice:
			switch t.Elem().Kind() {
			case reflect.Uint8:
				data := vv.Interface().([]byte)
				return int(binary.BigEndian.Uint32(data))
			default:
				fmt.Printf("Unsupported struct type %v", vv.Type().Name())
			}
		default:
			fmt.Printf("Unsupported type %v", vv.Type().Name())
		}
	}

	return
}

func Itf2Int64(val interface{}) (res int64) {
	if val == nil {
		return
	}

	if value, ok := val.(int64); ok {
		return value
	} else {
		t := reflect.TypeOf(val)
		vv := reflect.Indirect(reflect.ValueOf(val))
		switch t.Kind() {
		//checked
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return vv.Int()

			//checked
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int64(vv.Uint())

		//checked
		case reflect.Float32, reflect.Float64:
			return int64(vv.Float())

			//checked
		case reflect.String:
			return StrToInt64(vv.String())
			/*if i, err := strconv.ParseInt(vv.String(), 10, 0); err != nil {
				fmt.Printf("Unsupported type %v", vv.Type().Name())
			} else {
				return i
			}
			*/
		case reflect.Array, reflect.Slice:
			switch t.Elem().Kind() {
			case reflect.Uint8:
				data := vv.Interface().([]byte)
				return int64(binary.BigEndian.Uint64(data))
			default:
				fmt.Printf("Unsupported struct type Itf2Int64(%v)", vv.Type().Name())
			}
		default:
			fmt.Printf("Unsupported type Itf2Int64(%v)", vv.Type().Name())
		}
	}

	return
}

func Itf2Float(val interface{}) (res float64) {
	if val == nil {
		return
	}

	if value, ok := val.(float64); ok {
		return value
	} else {
		t := reflect.TypeOf(val)
		vv := reflect.Indirect(reflect.ValueOf(val))
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return vv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return vv.Float()
		case reflect.String:
			if f, err := strconv.ParseFloat(vv.String(), 64); err != nil {
				fmt.Printf("Unsupported type Itf2Float(%v) error : %s", vv.Type().Name(), err.Error())
			} else {
				return f
			}
		default:
			fmt.Printf("Unsupported type Itf2Float(%v)", vv.Type().Name())
		}
	}

	return
}

func Itf2Float32(val interface{}) (res float32) {
	if value, ok := val.(float32); ok {
		return value
	}

	return float32(Itf2Float(val))
}

func Itf2Str(val interface{}) (res string) {
	if val == nil {
		return
	}

	if value, ok := val.(string); ok {
		return value
	} else {
		t := reflect.TypeOf(val)
		vv := reflect.Indirect(reflect.ValueOf(val))
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return strconv.FormatInt(vv.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return strconv.FormatUint(vv.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			return strconv.FormatFloat(vv.Float(), 'f', -1, 64)
		case reflect.String:
			return vv.String()
		case reflect.Array, reflect.Slice:
			switch t.Elem().Kind() {
			case reflect.Uint8:
				data := vv.Interface().([]byte)
				return string(data)
			default:
				fmt.Printf("Unsupported struct type Itf2Str(%v)", vv.Type().Name())
			}
		//时间类型
		case reflect.Struct:
			var c_TIME_DEFAULT time.Time
			TimeType := reflect.TypeOf(c_TIME_DEFAULT)
			if t.ConvertibleTo(TimeType) {
				return vv.Convert(TimeType).Interface().(time.Time).Format(time.RFC3339Nano)
			} else {
				fmt.Printf("Unsupported struct type Itf2Str(%v)", vv.Type().Name())
			}
		case reflect.Bool:
			return strconv.FormatBool(vv.Bool())
		case reflect.Complex128, reflect.Complex64:
			return fmt.Sprintf("%v", vv.Complex())
		/* TODO: unsupported types below
		   case reflect.Map:
		   case reflect.Ptr:
		   case reflect.Uintptr:
		   case reflect.UnsafePointer:
		   case reflect.Chan, reflect.Func, reflect.Interface:
		*/
		default:
			fmt.Printf("Unsupported struct type Itf2Str(%v)", vv.Type().Name())
		}
	}

	return
}

func Itf2Time(val interface{}) (res time.Time) {
	if val == nil {
		return
	}

	if value, ok := val.(time.Time); ok {
		return value
	} else {
		t := reflect.TypeOf(val)
		vv := reflect.Indirect(reflect.ValueOf(val))
		//fmt.Println("datetime21", t, t.Kind(), val, vv, vv.String())
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return time.Unix(vv.Int(), 0)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return time.Unix(int64(vv.Uint()), 0)

		case reflect.String:

			if tm, err := time.Parse("2006-01-02 15:04:05", vv.String()); err != nil {
				fmt.Printf("Unsupported type Itf2Time(%v) error : %s", vv.Type().Name(), err.Error())
				//fmt.Println("String1:", val, vv.String(), err.Error())
			} else {
				fmt.Printf("Unsupported type Itf2Time(%v)", tm)
				return tm
			}
		case reflect.Struct:
			var c_TIME_DEFAULT time.Time
			TimeType := reflect.TypeOf(c_TIME_DEFAULT)
			//fmt.Println("datetime22", t, t.Kind(), t.ConvertibleTo(TimeType))

			if t.ConvertibleTo(TimeType) {
				return vv.Interface().(time.Time)
			} else {
				fmt.Printf("Unsupported struct type Itf2Time(%v)", vv.Type().Name())
			}
		default:
			fmt.Printf("Unsupported type Itf2Time(%v)", vv.Type().Name())
		}
	}

	return
}

// transfer a map or struct to interface map
func Itf2ItfMap(value interface{}) map[string]interface{} {
	if m, ok := value.(map[string]interface{}); ok {
		return m
	} else if m, ok := value.(map[string]string); ok {
		newMap := make(map[string]interface{})
		for k, v := range m {
			newMap[k] = v
		}
		return newMap
	} else {
		return Struct2ItfMap(value)
	}
}

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

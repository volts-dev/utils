package utils

import (
	"reflect"
)

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

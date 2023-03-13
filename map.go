package utils

import (
	"sync"
	"sync/atomic"
)

// "reflect"
//
//	"strings"
type (
	//TMap map[string]interface{}

	// TMap is a map with lock
	TMap struct {
		sync.RWMutex
		len  int32
		data map[interface{}]interface{}
	}
)

func MapToAnyList[T map[string]string | map[string]any](maps ...T) []any {
	result := make([]any, len(maps))
	for _, m := range maps {
		result = append(result, m)
	}

	return result
}

func NewMap() *TMap {
	return &TMap{
		data: make(map[interface{}]interface{}),
	}
}

func (self *TMap) Clear(k interface{}) {
	self.Lock()
	self.data = make(map[interface{}]interface{})
	self.Unlock()
}

// Get from maps return the k's value
func (self *TMap) Get(k interface{}) interface{} {
	self.RLock()
	defer self.RUnlock()
	if val, ok := self.data[k]; ok {
		return val
	}
	return nil
}

// Set Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (self *TMap) Set(k interface{}, v interface{}) bool {
	self.Lock()
	defer self.Unlock()
	self.data[k] = v
	atomic.StoreInt32(&self.len, self.len+1)

	return true
}

// Check Returns true if k is exist in the map.
func (self *TMap) IsExist(k interface{}) bool {
	self.Lock()
	defer self.Unlock()
	_, ok := self.data[k]
	return ok
}

// Delete the given key and value.
func (self *TMap) Delete(k interface{}) {
	self.Lock()
	delete(self.data, k)
	atomic.StoreInt32(&self.len, self.len-1)
	self.Unlock()
}

// Items returns all items in safemap.
func (self *TMap) Items() map[interface{}]interface{} {
	r := make(map[interface{}]interface{})
	self.RLock()
	for k, v := range self.data {
		r[k] = v
	}
	self.RUnlock()
	return r
}

func (self *TMap) Range(fn func(k, v interface{}) bool) {
	self.RLock()
	defer self.RUnlock()
	for k, v := range self.data {
		if !fn(k, v) {
			break
		}
	}

}

// Count returns the number of items within the map.
func (self *TMap) Count() int {
	return int(atomic.LoadInt32(&self.len))
}

/*
func InMap(key string) bool {
	if str, has := item[field]; has {
		result = append(result, str.(string))
	}
}
*/
// 获取[]Map的指定Field
func GetStrMapField(lst interface{}, field string) (result []string) {
	if lst == nil {
		return nil
	}
	var str string
	if m, ok := lst.([]map[string]interface{}); ok {
		for _, item := range m {
			if str, has := item[field]; has {
				result = append(result, str.(string))
			}
		}

	} else if m, ok := lst.([]map[string]string); ok {
		for _, item := range m {
			if str = item[field]; str != "" {
				result = append(result, str)
			}
		}
	} else if m, ok := lst.([]map[string][]byte); ok {
		for _, item := range m {
			if str = string(item[field]); str != "" {
				result = append(result, str)
			}
		}
	}

	return
}

func MergeMaps(to map[string]interface{}, from ...map[string]interface{}) map[string]interface{} {
	if from != nil { //如果有res
		if to == nil { //且没有des
			to = make(map[string]interface{})
		}

		//两者都有合并
		for i := 0; i < len(from); i++ {
			for key, value := range from[i] {
				to[key] = value
			}
		}
	}

	return to
}

func StrMap2ItfMap(m map[string]string) (res_map map[string]interface{}) {
	//res_map=make(map[string]interface{})
	return nil
}

func ItfMap2StrMap(m map[string]interface{}) map[string]string {
	return nil
}

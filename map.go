package utils

import (
	"sync"
)

//"reflect"
//	"strings"
type (
	//TMap map[string]interface{}

	// TMap is a map with lock
	TMap struct {
		lock *sync.RWMutex
		data map[interface{}]interface{}
	}
)

func NewMap() *TMap {
	return &TMap{
		lock: new(sync.RWMutex),
		data: make(map[interface{}]interface{}),
	}
}

// Get from maps return the k's value
func (self *TMap) Get(k interface{}) interface{} {
	self.lock.RLock()
	defer self.lock.RUnlock()
	if val, ok := self.data[k]; ok {
		return val
	}
	return nil
}

// Set Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (self *TMap) Set(k interface{}, v interface{}) bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	if val, ok := self.data[k]; !ok {
		self.data[k] = v
	} else if val != v {
		self.data[k] = v
	} else {
		return false
	}
	return true
}

// Check Returns true if k is exist in the map.
func (self *TMap) Check(k interface{}) bool {
	self.lock.RLock()
	defer self.lock.RUnlock()
	_, ok := self.data[k]
	return ok
}

// Delete the given key and value.
func (self *TMap) Delete(k interface{}) {
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.data, k)
}

// Items returns all items in safemap.
func (self *TMap) Items() map[interface{}]interface{} {
	self.lock.RLock()
	defer self.lock.RUnlock()
	r := make(map[interface{}]interface{})
	for k, v := range self.data {
		r[k] = v
	}
	return r
}

// Count returns the number of items within the map.
func (self *TMap) Count() int {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return len(self.data)
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

func MergeMaps(aSrc, aDes map[string]interface{}) map[string]interface{} {
	if aSrc == nil { //如果没有res
		return aDes
	} else { //如果有res
		if aDes == nil { //且没有des
			aDes = aSrc
			return aSrc
		} else { //两者都有合并
			for key, value := range aSrc {
				aDes[key] = value
			}

		}
	}
	return aDes
}

func StrMap2ItfMap(m map[string]string) (res_map map[string]interface{}) {
	//res_map=make(map[string]interface{})
	return nil
}

func ItfMap2StrMap(m map[string]interface{}) map[string]string {
	return nil
}

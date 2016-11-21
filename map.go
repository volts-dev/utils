package utils

//"reflect"
//	"strings"
type (
	TMap map[string]interface{}
)

func (self TMap) Contain(key string) (result bool) {
	_, result = self[key]
	return
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

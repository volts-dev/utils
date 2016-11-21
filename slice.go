package utils

func Insert(slice []interface{}, index int, value interface{}) {
	// Grow the slice by one element.
	// make([]Token, len(self.Child)+1)
	// self.Child[0 : len(self.Child)+1]
	slice = append(slice, value)
	// Use copy to move the upper part of the slice out of the way and open a hole.

	copy(slice[index+1:], slice[index:])
	// Store the new value.
	slice[index] = value
	// Return the result.
	//return slice
}

// TODO 改为范类型
func SlicInsert(slice []interface{}, index int, values ...interface{}) []interface{} {
	return append(slice[:index], append(values, slice[index:]...)...)
}

func SlicRemove(slice []interface{}, index int) []interface{} {
	return append(slice[:index], slice[index+1:]...)
}

func StringsCopy(slice []string) []string {
	tmp := make([]string, 0)
	return append(tmp, slice...)
}

// delete string from []string
func StringsDel(lst []string, key string) []string {
	for idx, str := range lst {
		if SameText(str, key) {
			lst = append(lst[:idx], lst[idx+1:]...)
		}
	}
	return lst
}

// 取交集
func StringsIntersect(a []string, b []string) (res []string) {
	res = make([]string, 0)
	for _, a_str := range a {
		for _, b_str := range b {
			if a_str == b_str {
				res = append(res, a_str)
			}
		}
	}
	return
}

// check if string in other strings
func InStrings(target string, other ...string) bool {
	if len(other) == 0 {
		return false
	}
	for _, str := range other {
		if target == str {
			return true
		}
	}
	return false
}

// 复制一个反转版
func Reversed(lst []string) (result []string) {
	result = make([]string, 0)
	for i := len(lst) - 1; i >= 0; i-- {
		result = append(result, lst[i])
	}
	return
}

func Strs2Itfs(m []string) (res_slice []interface{}) {
	res_slice = make([]interface{}, 0)

	for _, val := range m {
		res_slice = append(res_slice, val)
	}

	return
}

func Itfs2Strs(m []interface{}) (res []string) {
	res = make([]string, 0)

	for _, val := range m {
		res = append(res, Itf2Str(val))
	}

	return
}

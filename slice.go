package utils

import (
	"strings"
)

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
// return the index of the list otherwise -1 no match found
func InStrings(target string, other ...string) int {
	for idx, str := range other {
		if target == str {
			return idx
		}
	}
	return -1
}

func InInts(target int, other ...int) int {
	for idx, i := range other {
		if target == i {
			return idx
		}
	}
	return -1
}

func HasStrings(target string, other ...string) int {
	for _, str := range other {
		if idx := strings.Index(target, str); idx != -1 {
			return idx
		}
	}
	return -1
}

// 复制一个反转版
func Reversed(lst []string) (result []string) {
	result = make([]string, 0)
	for i := len(lst) - 1; i >= 0; i-- {
		result = append(result, lst[i])
	}
	return
}

func ReverseItfs(lst ...interface{}) (result []interface{}) {
	result = make([]interface{}, 0)
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

func IntsToStrs(m []int64) (res []string) {
	res = make([]string, 0)

	for _, val := range m {
		res = append(res, IntToStr(val))
	}

	return
}

func JoinQuote(list []string, quote, sep string) string {
	cnt := len(list)
	if cnt > 0 {
		n := (len(sep) * cnt) - 1
		for i := 0; i < len(list); i++ {
			n += len(list[i]) + 2
		}

		b := make([]byte, n)
		bp := 0 //copy(b, list[0])
		for i, s := range list {
			if i != 0 {
				bp += copy(b[bp:], sep)
			}
			bp += copy(b[bp:], quote)
			bp += copy(b[bp:], s)
			bp += copy(b[bp:], quote)
		}

		return string(b)
	}

	return ""
}

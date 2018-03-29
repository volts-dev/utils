package utils

import (
	"bytes"
	"reflect"
	"strconv"
	"unicode"
)

var (
	special_bytes = []byte(`.\+*?|[]{}^$`)
)

// 只解析'true','false'
func IsBoolStr(str string) bool {
	if IsIntStr(str) {
		return false
	}

	if _, err := strconv.ParseBool(str); err == nil {
		return true //	fmt.Printf("%T, %v\n", s, s)
	}
	return false
}

func IsBoolItf(itf interface{}, test_number ...bool) bool {
	switch itf.(type) {
	case bool:
		return true
	case string:
		str := itf.(string)
		if IsIntStr(str) {
			return false
		}

		if _, err := strconv.ParseBool(str); err == nil {
			return true //	fmt.Printf("%T, %v\n", s, s)
		}
	case int, int64:
		if len(test_number) > 0 {
			if i, ok := itf.(int); ok && (i == 1 || i == 0) {
				return true
			}
		}
	default:
	}

	return false
}

func IsIntStr(str string) bool {
	if _, err := strconv.ParseInt(str, 10, 0); err == nil {
		return true //	fmt.Printf("%T, %v\n", s, s)
	}
	return false
}

// Repeat returns a new string consisting of count copies of the string s.
func Repeat(s string, count int) (result []string) {
	result = make([]string, 0)
	for count > 0 {
		result = append(result, s)
		count--
	}
	return
}

// Equal is a helper for comparing value equality, following these rules:
//  - Values with equivalent types are compared with reflect.DeepEqual
//  - int, uint, and float values are compared without regard to the type width.
//    for example, Equal(int32(5), int64(5)) == true
//  - strings and byte slices are converted to strings before comparison.
//  - else, return false.
func Equal(a, b interface{}) bool {
	if reflect.TypeOf(a) == reflect.TypeOf(b) {
		return reflect.DeepEqual(a, b)
	}
	switch a.(type) {
	case int, int8, int16, int32, int64:
		switch b.(type) {
		case int, int8, int16, int32, int64:
			return reflect.ValueOf(a).Int() == reflect.ValueOf(b).Int()
		}
	case uint, uint8, uint16, uint32, uint64:
		switch b.(type) {
		case uint, uint8, uint16, uint32, uint64:
			return reflect.ValueOf(a).Uint() == reflect.ValueOf(b).Uint()
		}
	case float32, float64:
		switch b.(type) {
		case float32, float64:
			return reflect.ValueOf(a).Float() == reflect.ValueOf(b).Float()
		}
	case string:
		switch b.(type) {
		case []byte:
			return a.(string) == string(b.([]byte))
		}
	case []byte:
		switch b.(type) {
		case string:
			return b.(string) == string(a.([]byte))
		}
	}
	return false
}

func IsSpecialByte(ch byte) bool {
	return bytes.IndexByte(special_bytes, ch) > -1
}

func IsAlphaByte(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func IsDigitByte(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func IsAlnumByte(ch byte) bool {
	return IsAlphaByte(ch) || IsDigitByte(ch)
}

func IsAlphaNumericRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

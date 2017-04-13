package utils

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"strings"
)

func Md5(AStr string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(AStr)))
}

// 替代InStrings
func IdxOfStrings(target string, strs ...string) int {
	for idx, str := range strs {
		if target == str {
			return idx
		}
	}
	return -1
}

func Trim(s string) string {
	return strings.Trim(s, " ")
}

func SameText(AStrA string, AStrB string) bool {
	if strings.ToLower(AStrA) == strings.ToLower(AStrB) {
		return true
	} else {
		return false
	}

}

// convert like this: "HelloWorld" to "hello_world"
func SnakeCasedName(name string) string {
	newstr := make([]rune, 0)
	firstTime := true

	for _, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if firstTime == true {
				firstTime = false
			} else {
				newstr = append(newstr, '_')
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}

// # return a struct mame
func Obj2Name(obj interface{}) string {
	lName := reflect.Indirect(reflect.ValueOf(obj)).Type().Name()
	switch lName {
	case "string":
		{
			return obj.(string)
		}
	default:
		{
			return lName
		}
	}

	return ""
}

// convert like this: "HelloWorld" to "hello.world"
func DotCasedName(name string) string {
	newstr := make([]rune, 0)
	firstTime := true

	for _, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if firstTime == true {
				firstTime = false
			} else {
				newstr = append(newstr, '.')
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}

// convert like this: "hello_world" to "HelloWorld"
func TitleCasedName(name string) string {
	return _titleCasedName(name, false)
}

// convert like this: "hello_world" to "HelloWorld"
func TitleCasedNameWithSpace(name string) string {
	return _titleCasedName(name, true)
}
func _titleCasedName(name string, sapce bool) string {
	newstr := make([]rune, 0)
	upNextChar := true

	for _, chr := range name {
		switch {
		case upNextChar:
			upNextChar = false
			chr -= ('a' - 'A')
		case chr == '_':
			upNextChar = true
			if sapce {
				newstr = append(newstr, ' ')
			}
			continue
		}

		newstr = append(newstr, chr)
	}

	return string(newstr)
}

func PluralizeString(str string) string {
	if strings.HasSuffix(str, "y") {
		str = str[:len(str)-1] + "ie"
	}
	return str + "s"
}

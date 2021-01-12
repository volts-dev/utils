package utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
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

func TrimCasedName(name string) string {
	newstr := make([]rune, 0)

	for _, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			//chr -= ('A' - 'a')
			newstr = append(newstr, chr)
		}
	}

	return string(newstr)
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

// convert string from "hello_world" to "HelloWorld"
func TitleCasedName(name string) string {
	return _titleCasedName(name, false)
}

// convert string from "hello_world" to "HelloWorld"
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

// contains reports whether the string contains the byte c.
func contains(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

// Unquote interprets s as a single-quoted, double-quoted,
// or backquoted Go string literal, returning the string value
// that s quotes.  (If s is single-quoted, it would be a Go
// character literal; Unquote returns the corresponding
// one-character string.)
func Unquote(s string) (string, error) {
	n := len(s)
	if n < 2 {
		return "", errors.New("invalid quoted string")
	}
	quote := s[0]
	if quote != s[n-1] {
		return "", errors.New("lost the quote symbol on the end")
	}
	s = s[1 : n-1]

	if quote == '`' {
		if contains(s, '`') {
			return "", errors.New("the '`' symbol is found on the content")
		}
		return s, nil
	}

	if quote != '"' && quote != '\'' {
		return "", errors.New("lost the quote symbol on the begin")
	}

	//if contains(s, '\n') {
	//	//Println("contains(s, '\n')")
	//	return "contains(s, '\n')", strconv.ErrSyntax
	//}

	// Is it trivial?  Avoid allocation.
	if !contains(s, '\\') && !contains(s, quote) {
		switch quote {
		case '"':
			return s, nil
		case '\'':
			r, size := utf8.DecodeRuneInString(s)
			if size == len(s) && (r != utf8.RuneError || size != 1) {
				return s, nil
			}
		}
	}

	var runeTmp [utf8.UTFMax]byte
	buf := make([]byte, 0, 3*len(s)/2) // Try to avoid more allocations.
	for len(s) > 0 {
		c, multibyte, ss, err := strconv.UnquoteChar(s, quote)
		if err != nil {
			return "", err
		}
		s = ss
		if c < utf8.RuneSelf || !multibyte {
			buf = append(buf, byte(c))
		} else {
			n := utf8.EncodeRune(runeTmp[:], c)
			buf = append(buf, runeTmp[:n]...)
		}
		if quote == '\'' && len(s) != 0 {
			// single-quoted must be single character
			return "", strconv.ErrSyntax
		}
	}
	return string(buf), nil
}

package utils

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

type timeFormatType int

const (
	timeFormatNoTimezone timeFormatType = iota
	timeFormatNamedTimezone
	timeFormatNumericTimezone
	timeFormatNumericAndNamedTimezone
	timeFormatTimeOnly
)

type timeFormat struct {
	format string
	typ    timeFormatType
}

func (f timeFormat) hasTimezone() bool {
	// We don't include the formats with only named timezones, see
	// https://github.com/golang/go/issues/19694#issuecomment-289103522
	return f.typ >= timeFormatNumericTimezone && f.typ <= timeFormatNumericAndNamedTimezone
}

var (
	timeFormats = []timeFormat{
		{time.RFC3339, timeFormatNumericTimezone},
		{"2006-01-02T15:04:05", timeFormatNoTimezone}, // iso8601 without timezone
		{time.RFC1123Z, timeFormatNumericTimezone},
		{time.RFC1123, timeFormatNamedTimezone},
		{time.RFC822Z, timeFormatNumericTimezone},
		{time.RFC822, timeFormatNamedTimezone},
		{time.RFC850, timeFormatNamedTimezone},
		{"2006-01-02 15:04:05.999999999 -0700 MST", timeFormatNumericAndNamedTimezone}, // Time.String()
		{"2006-01-02T15:04:05-0700", timeFormatNumericTimezone},                        // RFC3339 without timezone hh:mm colon
		{"2006-01-02 15:04:05Z0700", timeFormatNumericTimezone},                        // RFC3339 without T or timezone hh:mm colon
		{"2006-01-02 15:04:05", timeFormatNoTimezone},
		{time.ANSIC, timeFormatNoTimezone},
		{time.UnixDate, timeFormatNamedTimezone},
		{time.RubyDate, timeFormatNumericTimezone},
		{"2006-01-02 15:04:05Z07:00", timeFormatNumericTimezone},
		{"2006-01-02", timeFormatNoTimezone},
		{"02 Jan 2006", timeFormatNoTimezone},
		{"2006-01-02 15:04:05 -07:00", timeFormatNumericTimezone},
		{"2006-01-02 15:04:05 -0700", timeFormatNumericTimezone},
		{time.Kitchen, timeFormatTimeOnly},
		{time.Stamp, timeFormatTimeOnly},
		{time.StampMilli, timeFormatTimeOnly},
		{time.StampMicro, timeFormatTimeOnly},
		{time.StampNano, timeFormatTimeOnly},
	}
)

// toInt returns the int value of v if v or v's underlying type
// is an int.
// Note that this will return false for int64 etc. types.
func toInt(v interface{}) (int, bool) {
	switch v := v.(type) {

	case time.Weekday:
		return int(v), true
	case time.Month:
		return int(v), true
	default:
		return 0, false
	}
}

func trimZeroDecimal(s string) string {
	var foundZero bool
	for i := len(s); i > 0; i-- {
		switch s[i-1] {
		case '.':
			if foundZero {
				return s[:i-1]
			}
		case '0':
			foundZero = true
		default:
			return s
		}
	}
	return s
}

func parseDateWith(s string, location *time.Location, formats []timeFormat) (d time.Time, e error) {
	for _, format := range formats {
		if d, e = time.Parse(format.format, s); e == nil {

			// Some time formats have a zone name, but no offset, so it gets
			// put in that zone name (not the default one passed in to us), but
			// without that zone's offset. So set the location manually.
			if format.typ <= timeFormatNamedTimezone {
				if location == nil {
					location = time.Local
				}
				year, month, day := d.Date()
				hour, min, sec := d.Clock()
				d = time.Date(year, month, day, hour, min, sec, d.Nanosecond(), location)
			}

			return
		}
	}
	return d, fmt.Errorf("unable to parse date: %s", s)
}

// StringToDateInDefaultLocation casts an empty interface to a time.Time,
// interpreting inputs without a timezone to be in the given location,
// or the local timezone if nil.
func stringToDateInDefaultLocation(s string, location *time.Location) (time.Time, error) {
	return parseDateWith(s, location, timeFormats)
}

// ToTimeInDefaultLocationE casts an empty interface to time.Time,
// interpreting inputs without a timezone to be in the given location,
// or the local timezone if nil.
func toTimeInDefaultLocation(i interface{}, location *time.Location) (tim time.Time, err error) {
	i = indirect(i)

	switch v := i.(type) {
	case time.Time:
		return v, nil
	case string:
		return stringToDateInDefaultLocation(v, location)
	case json.Number:
		s := ToInt64(v)
		//if err1 != nil {
		//	return time.Time{}, fmt.Errorf("unable to cast %#v of type %T to Time", i, i)
		//}
		return time.Unix(s, 0), nil
	case int:
		return time.Unix(int64(v), 0), nil
	case int64:
		return time.Unix(v, 0), nil
	case int32:
		return time.Unix(int64(v), 0), nil
	case uint:
		return time.Unix(int64(v), 0), nil
	case uint64:
		return time.Unix(int64(v), 0), nil
	case uint32:
		return time.Unix(int64(v), 0), nil
	default:
		return time.Time{}, fmt.Errorf("unable to cast %#v of type %T to Time", i, i)
	}
}

// ToTimeE casts an interface to a time.Time type.
func ToTime(i interface{}) (tim time.Time) {
	v, err := toTimeInDefaultLocation(i, time.UTC)
	if err != nil {
		fmt.Println(err)
	}
	return v
}

// ToIntE casts an interface to an int type.
func ToInt(i interface{}) int {
	i = indirect(i)

	switch v := i.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case int32:
		return int(v)
	case int16:
		return int(v)
	case int8:
		return int(v)
	case uint:
		return int(v)
	case uint64:
		return int(v)
	case uint32:
		return int(v)
	case uint16:
		return int(v)
	case uint8:
		return int(v)
	case float64:
		return int(v)
	case float32:
		return int(v)
	case string:
		s, err := strconv.ParseInt(trimZeroDecimal(v), 0, 0)
		if err == nil {
			return int(s)
		}
		fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
		return 0
	case json.Number:
		return ToInt(string(v))
	case bool:
		if v {
			return 1
		}
		return 0
	case nil:
		return 0
	default:
		intv, ok := toInt(i)
		if ok {
			return intv
		}

		fmt.Errorf("unable to cast %#v of type %T to int", i, i)
		return 0
	}
}

// ToInt64E casts an interface to an int64 type.
func ToInt64(i interface{}) int64 {
	i = indirect(i)

	switch s := i.(type) {
	case int:
		return int64(s)
	case int64:
		return s
	case int32:
		return int64(s)
	case int16:
		return int64(s)
	case int8:
		return int64(s)
	case uint:
		return int64(s)
	case uint64:
		return int64(s)
	case uint32:
		return int64(s)
	case uint16:
		return int64(s)
	case uint8:
		return int64(s)
	case float64:
		return int64(s)
	case float32:
		return int64(s)
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err == nil {
			return v
		}
		fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
		return 0
	case json.Number:
		return ToInt64(string(s))
	case bool:
		if s {
			return 1
		}
		return 0
	case nil:
		return 0
	default:
		intv, ok := toInt(i)
		if ok {
			return int64(intv)
		}

		fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
		return 0
	}
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

func SliceByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func JsonBodyAsMap(body []byte) (m map[string]interface{}, err error) {
	err = json.Unmarshal(body, &m)
	//	LogErr(err)
	return
}

func BoolToStr(b bool) (str string) {
	return strconv.FormatBool(b)
}

func IntToStr(i interface{}) string {
	switch i.(type) {
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", i)
	default:
		return "0"
	}
}

// int,int64 to bytes
func Int64ToBytes(val interface{}) (res []byte) {
	return big.NewInt(val.(int64)).Bytes()
}

/*
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}*/

func FloatToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func StrToFloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		//fmt.Errorf(err.Error())
		fmt.Printf("faild to convert StrToFloat(%s) with error : %s", str, err.Error())
	}

	return f
}

// HexToBytes converts a hex string representation of bytes to a byte representation
func HexToBytes(h string) []byte {
	s, err := hex.DecodeString(h)
	if err != nil {
		fmt.Printf("faild to convert BytesToFloat(%s) with error : %s", h, err.Error())
		return []byte("")
	}
	return s
}

// BytesToHex converts bytes to a hex string representation of bytes
func BytesToHex(b []byte) string {
	return hex.EncodeToString(b)
}

func BytesToInt64(buf []byte) int64 {
	//res, _ = strconv.ParseInt(string(val), 10, 64)

	return int64(binary.BigEndian.Uint64(buf))
}

func BytesToFloat(buf []byte) (res float64) {
	res, err := strconv.ParseFloat(string(buf), 32)
	if err != nil {
		fmt.Printf("faild to convert BytesToFloat(%s) with error : %s", string(buf), err.Error())
	}
	return
}

// Base64ToBytes converts from a b64 string to bytes
func Base64ToBytes(h string) []byte {
	s, err := base64.URLEncoding.DecodeString(h)
	if err != nil {
		fmt.Printf("faild to convert Base64ToBytes(%s) with error : %s", h, err.Error())
		return []byte("")
	}
	return s
}

// BytesToBase64 converts bytes to a base64 string representation
func BytesToBase64(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}

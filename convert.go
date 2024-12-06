package utils

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"math/big"
	"reflect"
	"strconv"
	"time"
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

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirectToStringerOrError returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil) or an implementation of fmt.Stringer
// or error,
func indirectToStringerOrError(a interface{}) interface{} {
	if a == nil {
		return nil
	}

	var errorType = reflect.TypeOf((*error)(nil)).Elem()
	var fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
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
		s, err := ToInt64E(v)
		if err != nil {
			return time.Time{}, err
		}
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
func ToTimeE(i interface{}) (time.Time, error) {
	return toTimeInDefaultLocation(i, time.UTC)
}

func ToTime(i interface{}) (tim time.Time) {
	v, _ := ToTimeE(i)
	return v
}

func ToBool(i any) bool {
	i = indirect(i)

	switch b := i.(type) {
	case bool:
		return b
	case nil:
		return false
	case int:
		return b != 0
	case int64:
		return b != 0
	case int32:
		return b != 0
	case int16:
		return b != 0
	case int8:
		return b != 0
	case uint:
		return b != 0
	case uint64:
		return b != 0
	case uint32:
		return b != 0
	case uint16:
		return b != 0
	case uint8:
		return b != 0
	case float64:
		return b != 0
	case float32:
		return b != 0
	case time.Duration:
		return b != 0
	case string:
		v, err := strconv.ParseBool(i.(string))
		if err != nil {
			return false
		}
		return v
	case json.Number:
		v := ToInt64(b)
		return v != 0 //, fmt.Errorf("unable to cast %#v of type %T to bool", i, i)
	default:
		return false //, fmt.Errorf("unable to cast %#v of type %T to bool", i, i)
	}
}

func ToInt(i interface{}) int {
	v, _ := ToIntE(i)
	return v
}

// ToIntE casts an interface to an int type.
func ToIntE(i interface{}) (int, error) {
	i = indirect(i)

	switch v := i.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case int32:
		return int(v), nil
	case int16:
		return int(v), nil
	case int8:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint64:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint8:
		return int(v), nil
	case float64:
		return int(v), nil
	case float32:
		return int(v), nil
	case string:
		if v != "" {
			s, err := strconv.ParseInt(trimZeroDecimal(v), 0, 0)
			if err == nil {
				return int(s), nil
			}
		}

		return 0, fmt.Errorf("unable to cast %#v of type %T to int", i, i)
	case json.Number:
		return ToIntE(string(v))
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		intv, ok := toInt(i)
		if ok {
			return intv, nil
		}

		return 0, fmt.Errorf("unable to cast %#v of type %T to int", i, i)
	}
}

func ToInt64(i interface{}) int64 {
	v, _ := ToInt64E(i)
	return v
}

// ToInt64E casts an interface to an int64 type.
func ToInt64E(i interface{}) (int64, error) {
	i = indirect(i)

	switch v := i.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case int32:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case string:
		if v != "" {
			v, err := strconv.ParseInt(trimZeroDecimal(v), 0, 0)
			if err == nil {
				return v, nil
			}
		}

		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	case json.Number:
		return ToInt64E(string(v))
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		intv, ok := toInt(i)
		if ok {
			return int64(intv), nil
		}

		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	}
}

// ToString casts an interface to a string type.
func ToString(i interface{}) string {
	v, _ := ToStringE(i)
	return v
}

// TODO 优化indirectToStringerOrError
// ToStringE casts an interface to a string type.
func ToStringE(i interface{}) (string, error) {

	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(s), 10), nil
	case json.Number:
		return s.String(), nil
	case []byte:
		return string(s), nil
	case template.HTML:
		return string(s), nil
	case template.URL:
		return string(s), nil
	case template.JS:
		return string(s), nil
	case template.CSS:
		return string(s), nil
	case template.HTMLAttr:
		return string(s), nil
	case nil:
		return "", nil
	}

	i = indirectToStringerOrError(i)
	switch s := i.(type) {
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		return "", fmt.Errorf("unable to cast %#v of type %T to string", i, i)
	}
}

// ToFloat32E casts an interface to a float32 type.
func ToFloat32(i interface{}) float32 {
	i = indirect(i)

	intv, ok := toInt(i)
	if ok {
		return float32(intv)
	}

	switch s := i.(type) {
	case float64:
		return float32(s)
	case float32:
		return s
	case int64:
		return float32(s)
	case int32:
		return float32(s)
	case int16:
		return float32(s)
	case int8:
		return float32(s)
	case uint:
		return float32(s)
	case uint64:
		return float32(s)
	case uint32:
		return float32(s)
	case uint16:
		return float32(s)
	case uint8:
		return float32(s)
	case string:
		v, err := strconv.ParseFloat(s, 32)
		if err == nil {
			return float32(v)
		}
		fmt.Printf("unable to cast %#v of type %T to float32", i, i)
		return 0
	case json.Number:
		v, err := s.Float64()
		if err == nil {
			return float32(v)
		}
		fmt.Printf("unable to cast %#v of type %T to float32", i, i)
		return 0
	case bool:
		if s {
			return 1
		}
		return 0
	case nil:
		return 0
	default:
		fmt.Printf("unable to cast %#v of type %T to float32", i, i)
		return 0
	}
}

// ToFloat64E casts an interface to a float64 type.
func ToFloat64(i interface{}) float64 {
	i = indirect(i)

	intv, ok := toInt(i)
	if ok {
		return float64(intv)
	}

	switch s := i.(type) {
	case float64:
		return s
	case float32:
		return float64(s)
	case int64:
		return float64(s)
	case int32:
		return float64(s)
	case int16:
		return float64(s)
	case int8:
		return float64(s)
	case uint:
		return float64(s)
	case uint64:
		return float64(s)
	case uint32:
		return float64(s)
	case uint16:
		return float64(s)
	case uint8:
		return float64(s)
	case string:
		v, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return v
		}
		fmt.Printf("unable to cast %#v of type %T to float64", i, i)
		return 0
	case json.Number:
		v, err := s.Float64()
		if err == nil {
			return v
		}
		fmt.Printf("unable to cast %#v of type %T to float64", i, i)
		return 0
	case bool:
		if s {
			return 1
		}
		return 0
	case nil:
		return 0
	default:
		fmt.Printf("unable to cast %#v of type %T to float64", i, i)
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

func JsonBodyAsMap(body []byte) (m map[string]interface{}, err error) {
	err = json.Unmarshal(body, &m)
	//	LogErr(err)
	return
}

// int,int64 to bytes
func Int64ToBytes(val interface{}) (res []byte) {
	return big.NewInt(val.(int64)).Bytes()
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

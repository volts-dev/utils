package utils

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"unsafe"
)

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

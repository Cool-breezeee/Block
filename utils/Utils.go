package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"unsafe"
)

// Int64ToBytes int64转字节数组
func Int64ToBytes(value int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, value)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// JsonToArray Json转成数组
func JsonToArray(jsonStr string) []string {
	var arr []string
	if err := json.Unmarshal([]byte(jsonStr), &arr); err != nil {
		fmt.Println(err)
		return nil
	}
	return arr
}

func ByteSliceToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

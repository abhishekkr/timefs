package fs

import (
	"log"
	"strconv"
)

func Int32ToStr(n int32) string {
	return strconv.Itoa(int(n))
}

func StrToInt32(s string) int32 {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Println("[fs.StrToInt32] failure to convert", s)
	}
	return int32(n)
}

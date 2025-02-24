package base62

import (
	"math"
	"strings"
)

// 62进制转换

// 0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ

const base62Str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Int2String 整数转62进制
func Int2String(seq uint64) string {
	if seq == 0 {
		return string(base62Str[0])
	}
	res := make([]byte, 0)
	for seq > 0 {
		t := seq % 62
		seq = seq / 62
		// todo 考虑性能问题
		res = append([]byte{base62Str[t]}, res...)
	}
	return string(res)
}

// String2Int 62进制转整数
func String2Int(str string) uint64 {
	bl := []byte(str)
	bl = reverse(bl)
	res := uint64(0)
	for idx, b := range bl {
		t := math.Pow(62, float64(idx))
		res += uint64(strings.IndexByte(base62Str, b)) * uint64(t)
	}
	return res
}

func reverse(s []byte) []byte {
	for i := 0; i < len(s)/2; i++ {
		s[i], s[len(s)-1-i] = s[len(s)-1-i], s[i]
	}
	return s
}
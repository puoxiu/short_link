package md5

import (
	"crypto/md5"
	"encoding/hex"
)

var (
	secret = []byte("xingxing_hello")
)

// todo
// 计算密码的MD5值
func EncryptPassword(password string) string {
	return password + string(secret)
}

// Md5 md5加密
func Cal(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	md5 := hash.Sum(nil)
	md5Str := hex.EncodeToString(md5)

	return md5Str
}
package md5

var (
	secret = []byte("xingxing_hello")
)

// todo
// 计算密码的MD5值
func EncryptPassword(password string) string {
	return password + string(secret)
}
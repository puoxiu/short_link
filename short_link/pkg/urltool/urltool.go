package urltool

import (
	"errors"
	"net/url"
	"path"
)

// GetBasePath returns the base path of a URL.
func GetBasePath(urlStr string) (string, error) {
	myUrl, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	if len(myUrl.Host) == 0 {
		return "", errors.New("url host is empty")
	}
    // 处理空路径和以斜杠结尾的路径
    if myUrl.Path == "" || myUrl.Path == "/" {
        return "", nil
    }
	return path.Base(myUrl.Path), nil
}
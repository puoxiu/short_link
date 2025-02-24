package connect

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var client = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
	Timeout: 2 * time.Second,
}

// Get 判断url是否可访问
func Get(url string) bool {
	resp, err := client.Get(url)

	if err != nil {
		logx.Errorw("http get failed", logx.LogField{Key: "err", Value: err.Error()})
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
package middleware

import (
	"context"
	"net"
	"net/http"
	"short_link_pro/sl_redict/redict_api/constants"
	"strings"
)

type ClientIPMiddleware struct {
}

func NewClientIPMiddleware() *ClientIPMiddleware {
	return &ClientIPMiddleware{}
}


func (m *ClientIPMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var User_IP, User_Agent string
		User_Agent = r.Header.Get("User-Agent")
		
		xffHeader := r.Header.Get("X-Forwarded-For")	
		if xffHeader != "" {
			ips := strings.Split(xffHeader, ",")
			User_IP = strings.TrimSpace(ips[0])
		} else {
            // 提取 IP 地址，去除端口信息
            host, _, err := net.SplitHostPort(r.RemoteAddr)
            if err != nil {
                User_IP = r.RemoteAddr
            } else {
                User_IP = host
            }
		}
		// todo: ip agent是否有可能为空
		
		reqCtx := r.Context()
		ctx := context.WithValue(reqCtx, constants.UserAgentKey, User_Agent)
		ctx = context.WithValue(ctx, constants.UserIPKey, User_IP)
		newReq := r.WithContext(ctx)

		next(w, newReq)
	}
}
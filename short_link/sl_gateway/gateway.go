package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"short_link_pro/common/etcd"
	gateway_error "short_link_pro/sl_gateway/error"
	"strings"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "settings.yaml", "the config file")

// 根据服务名称映射到对应的服务地址
// var ServiceMap = map[string]string{
// 	"auth":  "http://127.0.0.1:8961", // 认证服务
// 	"user":  "http://127.0.0.1:8962", // 用户服务
// 	// "chat":  "http://127.0.0.1:8963", // 聊天服务
// 	// "group": "http://127.0.0.1:8964", // 群组服务
// }

type Config struct {
	Addr string
	Etcd string
}

// auth 函数用于认证请求
func auth(authAddr string, res http.ResponseWriter, req *http.Request) (bool) {
	authReq, _ := http.NewRequest("POST", authAddr, nil)
	authHeader := req.Header.Get("Authorization")
	authReq.Header = req.Header
	authReq.Header.Set("ValidPath", req.URL.Path)
	authReq.Header.Set("Token", authHeader)
	authRes, err := http.DefaultClient.Do(authReq)	
	if err != nil {
		// 认证服务错误
		gateway_error.FailResponse("认证服务错误1", res)
		return false
	}
	defer authRes.Body.Close()
	type Response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data *struct {
			Username string `json:"username"`
		} `json:"data"`
	}
	var authResponse Response
	byteData,_ := io.ReadAll(authRes.Body)
	err = json.Unmarshal(byteData, &authResponse)
	if err != nil {
		// 认证服务错误
		fmt.Println(err)
		gateway_error.FailResponse("认证服务错误2", res)

		return false
	}

	if authResponse.Code != 0 {
		res.Write(byteData)
		return false
	}
	// 认证通过
	// 请求头中添加用户ID和角色信息--用于后续服务
	if authResponse.Data != nil {
		req.Header.Set("Username", authResponse.Data.Username)
	}

	return true
}

// proxy 函数用于代理请求
func proxy(serverice string, res http.ResponseWriter, req *http.Request) {
	addr := etcd.GetServiceAddress(config.Etcd, serverice + "_api")
	if addr == "" {
		gateway_error.FailResponse("服务错误", res)
		return
	}

	// 在 Go 的 HTTP 包中，req.Body 是一个 io.ReadCloser 类型的对象，它只能被读取一次
	var body []byte
	var err error
	if req.Body != nil {
		body, err = io.ReadAll(req.Body)
		if err != nil {
			logx.Error(err)
			gateway_error.FailResponse("服务错误", res)
			return
		}
		// 重置原始请求体以便后续使用（例如日志记录）
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}
	remoteAddr := strings.Split(req.RemoteAddr, ":")		// 获取客户端 IP 地址

	url := fmt.Sprintf("http://%s%s", addr, req.URL.String())
	// proxyReq, err := http.NewRequest(req.Method, url, req.Body)
	proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))
	if err != nil {
		logx.Error(err)
		gateway_error.FailResponse("服务错误", res)
		return
	}

	// 设置请求头
	proxyReq.Header = req.Header
	proxyReq.Header.Del("ValidPath")
	proxyReq.Header.Set("X-Forwarded-For", remoteAddr[0])	// 设置请求头 以便服务端获取客户端 IP 地址
	response, err := http.DefaultClient.Do(proxyReq)	
	if err != nil {
		fmt.Println(err)
		gateway_error.FailResponse("服务错误", res)
		return
	}
	defer response.Body.Close()

	io.Copy(res, response.Body)
}

// gateway 函数是主网关处理逻辑，用于将请求转发到相应的服务
func gateway(res http.ResponseWriter, req *http.Request) {
	fmt.Println("gateway")
	// 请求认证服务地址
	authAddr := etcd.GetServiceAddress(config.Etcd, "auth_api")
	authUrl := fmt.Sprintf("http://%s/api/auth/authentication", authAddr)
	// 认证--token鉴权
	if !auth(authUrl, res, req) {
		fmt.Println("认证失败")
		return
	}

	// 使用正则表达式提取请求路径中的服务名称
	regex, err := regexp.Compile(`/api/(.*?)/`)
	if err != nil {
		logx.Error(err)
		// 匹配失败--检查服务名称是否正确
		gateway_error.FailResponse("服务错误", res)
		return
	}
	list := regex.FindStringSubmatch(req.URL.Path)
	if len(list) != 2 {
		// 匹配失败--检查服务名称是否正确
		gateway_error.FailResponse("服务错误", res)
		return
	}

	serverice := list[1]
	proxy(serverice, res, req)
}

// TODO: 完善对应配置文件
// 注册处理函数，将所有请求都转发给 gateway 处理
var config Config


func main() {
	flag.Parse()
	conf.MustLoad(*configFile, &config)

	fmt.Println(config)

	http.HandleFunc("/", gateway)
	fmt.Printf("Starting server at %s...\n", config.Addr)

	http.ListenAndServe(config.Addr, nil)
}

package etcd

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/netx"
)

// DeliveryAddress 上送服务地址
func DeliveryAddress(etcdAddr string, serviceName string, addr string) {
	list := strings.Split(addr, ":")
	if len(list) != 2 {
		logx.Errorf("地址错误 %s", addr)
		return
	}

	if list[0] == "0.0.0.0" {
		ip := netx.InternalIp()
		addr = strings.ReplaceAll(addr, "0.0.0.0", ip)
	}
	client := InitEtcd(etcdAddr)
	_, err := client.Put(context.Background(), serviceName, addr)
	if err != nil {
		logx.Errorf("上送服务地址失败 %s", err)
	}
	logx.Infof("上送服务地址成功 %s %s", addr, serviceName)

}

// GetServiceAddress 获取服务地址
func GetServiceAddress(etcdAddr string, serviceName string) (addr string) {
	client := InitEtcd(etcdAddr)
	resp, err := client.Get(context.Background(), serviceName)
	if err == nil && len(resp.Kvs) > 0 {
		addr = string(resp.Kvs[0].Value)
	}

	return
}
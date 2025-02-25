package main

import (
	"flag"
	"fmt"

	"short_link_pro/common/etcd"
	"short_link_pro/sl_redict/redict_api/internal/config"
	"short_link_pro/sl_redict/redict_api/internal/handler"
	"short_link_pro/sl_redict/redict_api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/redict.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 将服务地址上送到 etcd： key：auth_api ---》 value：服务地址(ip:port)
	etcd.DeliveryAddress(c.Etcd, c.Name + "_api", fmt.Sprintf("%s:%d", c.Host, c.Port))

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

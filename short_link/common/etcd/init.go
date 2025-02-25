package etcd

import "github.com/coreos/etcd/clientv3"

// etcd客户端
func InitEtcd(addr string) *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{addr},
	})
	if err != nil {
		panic(err)
	}

	return cli
}
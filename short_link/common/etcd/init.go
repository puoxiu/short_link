
package etcd

import clientv3 "go.etcd.io/etcd/client/v3"


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
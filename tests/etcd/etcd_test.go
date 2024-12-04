package etcd

import (
	"context"
	"testing"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestEtcd(t *testing.T) {

	cli, err := clientv3.New(clientv3.Config{})
	if err != nil {
		return
	}

	cli.Put(context.TODO(), "key", "val")
	cli.Delete(context.TODO(), "key")
	cli.Get(context.TODO(), "key", clientv3.WithCreatedNotify())

}

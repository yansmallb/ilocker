package etcdclient

import (
	"fmt"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"strings"
	"time"
)

type Etcd struct {
	client client.KeysAPI
}

var ServcieTimeout = 60 * time.Second
var Path = "/yansmallb/machine/"

func NewEtcdClient(etcdpath string) (*Etcd, error) {
	endpoints := strings.Split(etcdpath, ",")
	cfg := client.Config{
		Endpoints: endpoints,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavctxailable
		HeaderTimeoutPerRequest: ServcieTimeout,
	}
	c, err := client.New(cfg)
	kapi := client.NewKeysAPI(c)

	etcd := new(Etcd)
	etcd.client = kapi
	return etcd, err
}

func (e *Etcd) CreateKey(key, value string, TTL time.Duration) error {
	cfg := &client.CreateInOrderOptions{
		TTL: TTL,
	}
	if _, err := e.client.CreateInOrder(context.Background(), Path+key, value, cfg); err != nil {
		return err
	}
	return nil
}

func (e *Etcd) GetKey(key string) (string, error) {
	response, err := e.client.Get(context.Background(), Path+key, nil)
	if err != nil {
		return "", err
	}
	return response.Node.Value, nil
}

func (e *Etcd) ListKey() ([]string, error) {
	cfg := &client.GetOptions{
		Recursive: true,
		Sort:      false,
		Quorum:    false,
	}
	response, err := e.client.Get(context.Background(), Path, cfg)
	if err != nil {
		return nil, err
	}
	fmt.Println(response)
	return nil, nil
}

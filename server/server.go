package server

import (
	"github.com/yansmallb/ilocker/etcdclient"
	"time"
)

var machineList []string

func Server(etcdpath string, hb time.Duration) error {
	e, err := etcdclient.NewEtcdClient(etcdpath)
	if err != nil {
		return err
	}
	for {
		machineList, err = e.ListKey()
		if err != nil {
			return err
		}
		time.Sleep(hb)
	}
	return nil
}

func GetMachineList() []string {
	return machineList
}

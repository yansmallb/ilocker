package join

import (
	"encoding/json"
	"github.com/yansmallb/ilocker/etcdclient"
	"github.com/yansmallb/ilocker/plugin"
	"net/http"
	"time"
)

type Machine struct {
	Info map[string]string
}

func Join(etcdpath string, hb time.Duration) error {
	// get local machine md5
	m, _ := plugin.GetMachineInfo()
	mathineMd5, err := plugin.GetMachineMd5(m)
	if err != nil {
		return err
	}
	ip, err := plugin.GetMachineIp()
	if err != nil {
		return err
	}
	e, err := etcdclient.NewEtcdClient(etcdpath)
	if err != nil {
		return err
	}
	for {
		e.CreateKey(ip, mathineMd5, hb)
		time.Sleep(hb)
	}
}

func ListenAndServe() error {
	http.HandleFunc("/machine", getMachineInfo)
	http.ListenAndServe(":2374", nil)
	return nil
}

func getMachineInfo(w http.ResponseWriter, r *http.Request) {
	m, _ := plugin.GetMachineInfo()
	machine := &Machine{
		Info: m,
	}
	json.NewEncoder(w).Encode(machine)
}

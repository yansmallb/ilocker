package main

import (
	"encoding/json"
	"fmt"
	"github.com/yansmallb/ilocker/plugin"
	"io/ioutil"
	"net/http"
	"os"
)

type Machine struct {
	Info map[string]string
}

type Machines struct {
	List []string
}

var command = os.Getenv("command")

func main() {
	if len(os.Args) < 3 {
		os.Exit(-1)
	}
	managePath := os.Args[1]
	joinPath := os.Args[2]
	host, err := GetMachine(joinPath)
	if err != nil {
		os.Exit(-1)
	}
	machines := GetMachineList(managePath)
	hostMd5, _ := plugin.GetMachineMd5(host.Info)
	for _, machine := range machines {
		if machine == hostMd5 {
			plugin.ExecShell(command)
			os.Exit(0)
		}
	}
	os.Exit(-1)
}

func GetMachine(joinPath string) (*Machine, error) {
	path := joinPath + "/machine"
	body, err := httpGet(path)
	if err != nil {
		fmt.Println(err)
	}
	machine := &Machine{}
	err = json.Unmarshal(body, machine)
	return machine, err
}

func GetMachineList(managePath string) []string {
	path := managePath + "/machines"
	body, err := httpGet(path)
	if err != nil {
		fmt.Println(err)
	}
	machines := &Machines{}
	json.Unmarshal(body, machines)
	return machines.List
}

func httpGet(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

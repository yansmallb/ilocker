package plugin

import (
	"crypto/md5"
	"encoding/hex"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
)

var ipTemplate = "192\\.168\\.(-d)*\\.(-d)*"

func GetMachineInfo() (map[string]string, error) {
	m := make(map[string]string, 10)
	m["arch"] = runtime.GOARCH
	m["os"] = runtime.GOOS
	m["hostname"], _ = os.Hostname()
	m["lspci"] = ExecShell("lspci")
	m["iostat"] = ExecShell("iostat ALL | wc -l")
	m["dmidecode"] = ExecShell("dmidecode | grep -i 'serial number'")
	m["hdparm"] = ExecShell("hdparm -I /dev/sda |grep -i 'serial number' ")
	return m, nil
}

func GetMachineMd5(m map[string]string) (string, error) {
	var machineInfo string
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		machineInfo += m[k]
	}

	decode := md5.New()
	decode.Write([]byte(machineInfo))                 // 需要加密的字符串为 sharejs.com
	machineMd5 := hex.EncodeToString(decode.Sum(nil)) // 输出加密结果

	return machineMd5, nil
}

func GetMachineIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if isMatch, _ := regexp.MatchString(ipTemplate, ipnet.IP.String()); isMatch {
					return ipnet.IP.String(), nil
				}
			}
		}
	}
	return "", nil
}

func ExecShell(shell string) string {
	rtn, _ := exec.Command("/bin/bash", shell).Output()
	return string(rtn)
}

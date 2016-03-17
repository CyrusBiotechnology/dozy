package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func getVersion() string {
	versionStr := fmt.Sprintf("v%v", VERSION[0])
	for i := 1; i < len(VERSION); i++ {
		versionStr = fmt.Sprintf("%v.%v", versionStr, VERSION[i])
	}
	return versionStr
}

func getUptime() (time.Duration, error) {
	uptime := time.Duration(0)
	uptime_str, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return uptime, err
	}
	uptime_seconds, err := strconv.Atoi(strings.Split(string(uptime_str), ".")[0])
	if err != nil {
		return uptime, err
	}
	uptime = time.Duration(int(uptime_seconds)) * time.Second
	return uptime, nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

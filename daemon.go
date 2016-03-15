// +build darwin linux
package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func daemon(minUptime time.Duration, keyPollInterval time.Duration) {
	uptime_str, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		panic(err)
	}
	uptime_seconds, err := strconv.Atoi(strings.Split(string(uptime_str), ".")[0])
	if err != nil {
		panic(err)
	}
	uptime := time.Duration(int(uptime_seconds)) * time.Second

	if uptime < minUptime {
		wait := (minUptime - uptime)
		Info.Println(fmt.Sprintf("waiting %v for configured minimum uptime", wait))
		time.Sleep(wait)
	}

	keysTicker := time.NewTicker(keyPollInterval)
	for _ = range keysTicker.C {
		fmt.Println("tick!")
	}
}

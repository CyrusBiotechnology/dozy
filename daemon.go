// +build darwin linux
package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type DaemonConf struct {
	// Omit to use fs watcher (not yet supported)
	KeyPollInterval time.Duration
	// Use etcd to record and receive census information
	EtcdCensusServer string
	// Duration after which to activate
	MinUptime time.Duration
}

func daemon(conf DaemonConf) {
	uptime_str, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		exitE(err)
	}
	uptime_secnds, err := strconv.Atoi(strings.Split(string(uptime_str), ".")[0])
	if err != nil {
		exitE(err)
	}
	uptime := time.Duration(int(uptime_secnds)) * time.Second

	if uptime < conf.MinUptime {
		wait := (conf.MinUptime - uptime)
		Info.Println(fmt.Sprintf("waiting %v for configured minimum uptime", wait))
		time.Sleep(wait)
	}

	keysTicker := time.NewTicker(conf.KeyPollInterval)
	keysTicker := time.NewTicker(conf.KeyPollInterval)
}

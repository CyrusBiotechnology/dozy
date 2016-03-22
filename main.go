// +build darwin linux
package main

import (
	"fmt"
)

var VERSION = [...]int{1, 1, 1}

func main() {
	logging()
	configure()
	Info.Println("minimum uptime:", settings.MinUptime, "locks valid for", settings.LockAge, "locks:", settings.Locks)

	if settings.DaemonMode {
		Info.Println("running as daemon")
		daemon(settings.MinUptime, settings.MaxUptime, settings.Daemon.MinPeers, settings.Daemon.MaxPeers, settings.Daemon.KeyPollInterval)
	} else {
		if locksPlaceExists, _ := exists(settings.Locks); !locksPlaceExists {
			panic("specified locks location doesn't exist")
		}

		valid, invalid, err := getLocks(settings.Locks, settings.LockAge)
		if err != nil {
			panic(err)
		}
		if len(invalid) > 0 {
			Error.Println(len(invalid), "stale locks found")
		}
		if len(valid) > 0 {
			panic(fmt.Sprintf("%v valid lock(s) found", len(valid)))
		}
	}

	running, err := getRunningContainers()
	if len(running) > 0 {
		Info.Println("no locks found, stopping", len(running), "containers")
		stopAllRunningContainersWithRetry()
		if err != nil {
			panic(err)
		}
	} else {
		Info.Println("didn't find any running containers")
	}
	Info.Println("ready to doze :)")
}

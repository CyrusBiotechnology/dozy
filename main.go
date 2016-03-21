// +build darwin linux
package main

import (
	"fmt"
)

var VERSION = [...]int{1, 1, 1}

func main() {
	configure()
	logging()
	Info.Println(fmt.Sprintf("minimum uptime: %v, locks valid for %v, lock: %v",
		settings.MinUptime, settings.LockAge, settings.Locks))

	if settings.DaemonMode {
		Info.Println("running as daemon")
		daemon(settings.MinUptime, settings.Daemon.KeyPollInterval)
	} else {
		if locksPlaceExists, _ := exists(settings.Locks); !locksPlaceExists {
			panic("specified locks location doesn't exist")
		}

		valid, invalid, err := getLocks(settings.Locks, settings.LockAge)
		if len(invalid) > 0 {
			Error.Println(fmt.Printf("%v stale locks found", len(invalid)))
		}
		if len(valid) > 0 {
			panic(fmt.Sprintf("%v valid lock(s) found", len(valid)))
		}
		running, err := getRunningContainers()
		if len(running) > 0 {
			Info.Println(fmt.Printf("no locks found, stopping %v containers", len(running)))
			stopAllRunningContainersWithRetry()
			if err != nil {
				panic(err)
			}
		} else {
			Info.Println("didn't find any running containers")
		}
		Info.Println("ready to doze :)")
	}
}

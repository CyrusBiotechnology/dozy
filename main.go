// +build darwin linux
package main

import (
	"fmt"
)

var VERSION = [...]int{1, 0, 6}

func main() {
	configure()
	logging(fmt.Sprint(" `( ◔ ౪◔)´  dozy ", getVersion()))
	Info.Println(fmt.Sprintf("minimum uptime: %v, locks valid for %v, lock: %v",
		settings.MinUptime, settings.LockAge, settings.Locks))

	if locksPlaceExists, _ := exists(settings.Locks); !locksPlaceExists {
		panic("specified locks location doesn't exist")
	}

	if settings.DaemonMode {
		Info.Println("running as daemon")
		daemon(settings.MinUptime, settings.Daemon.KeyPollInterval)
	} else {
		valid, err := getValidLocks(settings.Locks, settings.LockAge)
		if len(valid) > 0 {
			panic("valid lock(s) found")
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

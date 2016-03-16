// +build darwin linux
package main

import (
	"fmt"
)

var VERSION = [...]int{1, 0, 6}

func main() {
	configure()
	logging(settings.Logs, fmt.Sprint(" `( ◔ ౪◔)´  dozy ", getVersion()))
	Info.Println(fmt.Sprintf("minimum uptime: %v, locks valid for %v, lock: %v",
		settings.MinUptime, settings.LockAge, settings.Locks))

	if locksPlaceExists, _ := exists(settings.Locks); !locksPlaceExists {
		panic("specified locks location doesn't exist")
	}

	if settings.Daemon {
		Info.Println("running as daemon")
		daemon(settings.MinUptime, settings.Daemon.KeyPollInterval)
	} else {
		valid, err := getValidLocks(settings.Locks, settings.LockAge)
		if err != nil {
			panic(err)
		}
		if len(valid) > 0 {
			panic("valid lock(s) found")
		}
	}
}

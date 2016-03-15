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

	getValidLocks(settings.Locks, settings.LockAge)
}

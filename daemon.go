// +build darwin linux
package main

import (
	"fmt"
	"time"
)

func daemon(minUptime time.Duration, keyPollInterval time.Duration) {
	// Should we process errors here?
	uptime, _ := getUptime()

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

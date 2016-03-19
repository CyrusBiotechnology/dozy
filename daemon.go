// +build darwin linux
package main

import (
	"fmt"
	"net"
	"time"
)

func daemon(minUptime time.Duration, keyPollInterval time.Duration) {
	// Should we process errors here?
	uptime, _ := getUptime()

	done := make(chan struct{})
	addr := net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 30000,
	}
	go censusServer(done, "udp4", &addr)

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

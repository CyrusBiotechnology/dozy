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
	listen := net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 19091,
	}
	bcast := net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 19091,
	}

	go Serve(done, "udp4", &listen, &bcast)

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

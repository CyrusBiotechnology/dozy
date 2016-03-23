// +build darwin linux
package main

import (
	"fmt"
	"github.com/CyrusBiotechnology/censusd"
	"net"
	"time"
)

func daemon(minUptime time.Duration, maxUptime time.Duration, minPeers int, maxPeers int, keyPollInterval time.Duration) error {
	// Should we process errors here?
	uptime, _ := getUptime()

	done := make(chan struct{})
	listen := net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 19091,
	}

	stats, err := censusd.Serve(done, "udp4", &listen)
	if err != nil {
		panic(err)
	}

	if uptime < minUptime {
		wait := (minUptime - uptime)
		Info.Println(fmt.Sprintf("waiting %v for configured minimum uptime", wait))
		time.Sleep(wait)
	}

	keysTicker := time.NewTicker(keyPollInterval)
	for _ = range keysTicker.C {
		uptime, _ := getUptime()
		stats.Mutex.RLock()
		peers := stats.Nodes
		stats.Mutex.RUnlock()
		if maxUptime > time.Duration(0) && uptime > maxUptime {
			Info.Println("this node is too old!", uptime)
			Info.Println("uptime:", uptime, "peers:", peers)
			close(done)
			return nil
		}
		if minUptime > uptime {
			continue
		}
		if peers <= minPeers {
			continue
		}
		Info.Println("conditions satisfied. uptime:", uptime, "peers:", peers, "/", minPeers)
		close(done)
		return nil
	}
	panic("how did we get here?")
}

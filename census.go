// +build darwin linux
package main

import (
	"net/http"
	"time"
)

// Call put() in a loop with backoff
func Etcd(interval time.Duration, server string) {
	t := time.NewTicker(interval)

	select {
	case <-t:
		http.NewRequest("PUT", server, nil)
	}
}

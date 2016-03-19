package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Node struct {
	LastContact time.Time // Timestamp @last time we heard from this peer
	Mutex       sync.RWMutex
}

type NodeList struct {
	Nodes     map[string]Node // string uniquely identifies a peer in this swarm
	MaxLength int             // How many peers do we track?
	Mutex     sync.RWMutex
}

func newNodeList(maxLength int) NodeList {
	return NodeList{
		Nodes:     map[string]Node{},
		MaxLength: maxLength,
		Mutex:     sync.RWMutex{},
	}
}

// node ids is a hash map with a pointer to an element in a linked list
// node

func newNode(nodeList *NodeList, node Node) {
	nodeList.Mutex.Lock()
	defer nodeList.Mutex.Unlock()
	nodeList.Nodes[]
}

// Census server listens for UDP broadcast packets and estimates peer
// population. Note that this may not work across network boundaries depending
// on firewall setup.
func censusServer(done <-chan struct{}, proto string, addr *net.UDPAddr) {
	peers := newNodeList(1000)
	cullTicker := time.NewTicker(time.Second)
	socket, err := net.ListenUDP(proto, addr)
	if err != nil {
		return
	}
	go func() {
		for {
			select {
			case <-done:
				return
			case <-cullTicker.C:

			default:
				buf := make([]byte, 4096)
				n, remoteAddr, err := socket.ReadFromUDP(buf)
				if err != nil {
					return
				}
				fmt.Println("Got:", string(buf[0:n]), "From:", remoteAddr)
			}
		}
	}()
}

func getPeers() int {
	censusMutex.RLock()
	defer censusMutex.Unlock()
	population := census
	return population
}

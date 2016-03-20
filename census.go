package main

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"net"
	"regexp"
	"sync"
	"text/template"
	"time"
)

var uidRegex = regexp.MustCompile("^[A-Za-z0-9]{32}:")
var formatErr = errors.New("First 32 bytes of a message should contain a UID, followed by a colon")

type Server struct {
	Address *net.UDPAddr // Listen and broadcast
	UID     string       // Identify us to other peers in the swarm
}

type Beacon struct {
	Time   time.Time
	Sender string
}

type NodeGraph struct {
	Nodes   map[string]*list.Element // UID: pointer to item in history
	History *list.List               // Stores the time and sender of messages received in order
	Mutex   sync.RWMutex             // Protects against races
}

func NewNodeGraph() NodeGraph {
	return NodeGraph{
		Nodes:   make(map[string]*list.Element),
		History: list.New(),
		Mutex:   sync.RWMutex{},
	}
}

// Update node adds or updates a node.
func (ng *NodeGraph) updateNode(nodeUID string) {
	ng.Mutex.Lock()
	defer ng.Mutex.Unlock()
	ng.History.Remove(ng.Nodes[nodeUID])
	ng.Nodes[nodeUID] = ng.History.PushFront(Beacon{
		Time:   time.Now(),
		Sender: nodeUID,
	})
}

// Takes a message and runs the appropriate action.
func (ng *NodeGraph) processMessage(packet string) error {
	if uidRegex.Match([]byte(packet[0:33])) {
		ng.updateNode(packet[0:32])
		return nil
	} else {
		return formatErr
	}
}

func (ng *NodeGraph) calcInterval() time.Duration {
	ng.Mutex.Lock()
	defer ng.Mutex.Unlock()
	return time.Second * time.Duration(len(ng.Nodes))
}

func (ng *NodeGraph) gc() {
	ng.Mutex.Lock()
	defer ng.Mutex.Unlock()
	event := ng.History.Back().Value.(Beacon)
	threshold := time.Now().Add(-time.Duration(ng.History.Len()) * time.Second)
	for event.Time.Before(threshold) {
		delete(ng.Nodes, event.Sender)
		ng.History.Remove(ng.History.Back())
		event = ng.History.Back().Value.(Beacon)
	}
}

func (sc *Server) send(message string) error {
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   sc.Address.IP,
		Port: sc.Address.Port,
	})
	defer socket.Close()

	if err != nil {
		return err
	}
	socket.Write([]byte(message))
	return nil
}

// Keep the swarm notified of our existence.
func (sc *Server) doBeacon(stop <-chan struct{}, ng *NodeGraph) error {
	msgTmpl, err := template.New("message").Parse("{{.UID}}:")
	buf := new(bytes.Buffer)
	if err != nil {
		return err
	}
	for {
		after := time.After(time.Duration(ng.calcInterval()) * time.Second)
		select {
		case <-stop:
			return nil
		case <-after:
			err := msgTmpl.Execute(buf, sc)
			if err != nil {
				Error.Println("error executing template")
				continue
			}
			sc.send(buf.String())
			after = time.After(time.Duration(ng.calcInterval()) * time.Second)
		}
	}
}

// Census server listens for UDP broadcast packets and updates graph data.
// Note that this may not work "out of the box" across network boundaries
// depending on program and networking configuration.
func Serve(exit <-chan struct{}, proto string, address *net.UDPAddr) {
	uid, err := SecureRandomAlphaString(32)
	if err != nil {
		panic("could not generate uid!")
	}
	server := Server{
		Address: address,
		UID:     uid,
	}
	graph := NewNodeGraph()
	gc := time.NewTicker(time.Second)
	socket, err := net.ListenUDP(proto, address)
	if err != nil {
		return
	}
	server.doBeacon(exit, &graph)

	go func() {
		for {
			select {
			case <-exit:
				return
			case <-gc.C:
				graph.gc()
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

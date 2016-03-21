package main

import (
	"bytes"
	"container/list"
	"errors"
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

func (ng *NodeGraph) hasNode(UID string) bool {
	if _, ok := ng.Nodes[UID]; ok {
		return true
	} else {
		return false
	}
}

// Update node adds or updates a node.
func (ng *NodeGraph) updateNode(nodeUID string) {
	ng.Mutex.Lock()
	defer ng.Mutex.Unlock()
	if ng.hasNode(nodeUID) {
		ng.History.Remove(ng.Nodes[nodeUID])
	} else {
		Info.Println("new node:", nodeUID)
	}
	ng.Nodes[nodeUID] = ng.History.PushFront(Beacon{
		Time:   time.Now(),
		Sender: nodeUID,
	})
}

// Takes a message and runs the appropriate action.
func (ng *NodeGraph) processMessage(buffer []byte) error {
	if uidRegex.Match(buffer[0:33]) {
		ng.updateNode(string(buffer[0:32]))
		return nil
	} else {
		return formatErr
	}
}

func (ng *NodeGraph) calcInterval() time.Duration {
	ng.Mutex.Lock()
	defer ng.Mutex.Unlock()
	return time.Second * time.Duration(len(ng.Nodes)+1)
}

func (ng *NodeGraph) gc() {
	ng.Mutex.Lock()
	defer ng.Mutex.Unlock()
	if ng.History.Len() > 0 {
		event := ng.History.Back().Value.(Beacon)
		threshold := time.Now().Add(-time.Duration(ng.History.Len()+5) * time.Second)
		for event.Time.Before(threshold) {
			delete(ng.Nodes, event.Sender)
			ng.History.Remove(ng.History.Back())
			event = ng.History.Back().Value.(Beacon)
		}
	}
}

type Client struct {
	Broadcast *net.UDPAddr // Where to sent packets to
	UID       string       // Identify ourselves to peers
}

func (cl *Client) send(message string) error {
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   cl.Broadcast.IP,
		Port: cl.Broadcast.Port,
	})
	defer socket.Close()

	if err != nil {
		return err
	}
	socket.Write([]byte(message))
	return nil
}

// Keep the swarm notified of our existence.
func (cl *Client) doBeacon(stop <-chan struct{}, ng *NodeGraph) error {
	msgTmpl, err := template.New("message").Parse("{{.UID}}:")
	buf := new(bytes.Buffer)
	err = msgTmpl.Execute(buf, cl)
	if err != nil {
		Error.Println("error executing template")
		return err
	}
	msg := buf.String()
	if err != nil {
		return err
	}
	for {
		after := time.After(ng.calcInterval())
		select {
		case <-stop:
			return nil
		case <-after:
			cl.send(msg)
			after = time.After(ng.calcInterval())
		}
	}
}

// Census server listens for UDP broadcast packets and updates graph data.
// Note that this may not work "out of the box" across network boundaries
// depending on program and networking configuration.
func Serve(exit <-chan struct{}, proto string, listen *net.UDPAddr, bcast *net.UDPAddr) {
	uid, err := SecureRandomAlphaString(32)
	if err != nil {
		panic("could not generate uid!")
	}
	Info.Println("uid:", uid)
	client := Client{
		Broadcast: bcast,
		UID:       uid,
	}
	graph := NewNodeGraph()
	gc := time.NewTicker(time.Second)
	socket, err := net.ListenUDP(proto, listen)
	if err != nil {
		panic(err)
	}
	Info.Println("listening at:", listen.IP, listen.Port)
	go client.doBeacon(exit, &graph)

	go func() {
		for {
			select {
			case <-exit:
				return
			case <-gc.C:
				graph.gc()
			default:
				buf := make([]byte, 4096)
				_, _, err := socket.ReadFromUDP(buf)
				if err != nil {
					Info.Println("problem reading packet")
				}
				graph.processMessage(buf)
			}
		}
	}()
}

package census

import (
	//"net/http"
	//"time"

	"github.com/coreos/etcd/client"
	//"log"
)

type EtcdCensus struct {
	EtcdConfig client.Config
	Key        string
	Folder     string
}

//func (p *EtcdPollerConfig) New() {
//	c, err := client.New(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//	kapi := client.NewKeysAPI(c)
//	// set "/foo" key with "bar" value
//}
//
//func (p *EtcdPoller) Get() {
//	//
//}
//
//// Call put() in a loop with backoff. In this state
//func (p *EtcdPoller) Etcd(conf client.Config, server string) {
//	select {
//	case <-t:
//		http.NewRequest("PUT", server, nil)
//	}
//}

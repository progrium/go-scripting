package main

import (
	"log"
	"time"
	"fmt"

	"github.com/progrium/go-extensions"
	"github.com/progrium/duplex/prototype"
)

var observers = extensions.ExtensionPoint(new(EventObserver))

type EventObserver interface {
	Event(payload string)
}

type RemoteProxy struct{
	peer *duplex.Peer
	service string
	peerName string
}

type EventArgs struct {
	Payload string
}

func (p *RemoteProxy) Event(payload string) {
	// TODO: replace with "CallTo" (Call that uses OpenWith) using p.peerName
	err := p.peer.Call(p.service+".Event", &EventArgs{payload}, new(struct{}))
	if err != nil {
		log.Fatal(err)
	}
}


type Extensions struct {
	peer *duplex.Peer
}

type RegisterArgs struct {
	Interface string
	Name string
	Service string
}

type RegisterReply struct {}


func (e *Extensions) Register(args RegisterArgs, reply *RegisterReply) error {
	fmt.Println("register:", args.Interface, args.Name, args.Service)
	/* An attempt to do a dynamic proxy fails because of the types involved in doing Calls
	proxy := extensions.NewProxy(args.Interface, 
		func (method string, a []interface{}) interface{} {
			err := e.peer.Call(args.Service+"."+method, &EventArgs{a[0].(string)}, new(struct{}))
			if err != nil {
				log.Fatal(err)
			}
			return nil
		},
	) */
	extensions.RegisterWithName(args.Interface, &RemoteProxy{e.peer, args.Service, args.Name}, args.Name)
	return nil
}



func main() {
	peer := duplex.NewPeer()
	defer peer.Close()
	if err := peer.Bind("127.0.0.1:9877"); err != nil {
		log.Fatal(err)
	}

	peer.Register(&Extensions{peer})
	go peer.Serve()

	for {
		time.Sleep(5 * time.Second)
		for _, observer := range observers.All() {
			observer.(EventObserver).Event("Hello")
		}
	}
}
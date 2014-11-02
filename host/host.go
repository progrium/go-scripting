package main

import (
	"log"
	"fmt"

	"github.com/progrium/duplex/prototype"
)

type EventArgs struct {
	Payload string
}

type EventReply struct {}

type RegisterArgs struct {
	Interface string
	Name string
	Service string
}

type RegisterReply struct {}

type EventPrinter int

func (p *EventPrinter) Event(args EventArgs, reply *EventReply) error {
	fmt.Println("EventPrinter:", args.Payload)
	return nil
}

func main() {
	peer := duplex.NewPeer()
	defer peer.Close()
	if err := peer.Connect("127.0.0.1:9877"); err != nil {
		log.Fatal(err)
	}

	err := peer.Call("Extensions.Register", &RegisterArgs{"EventObserver", peer.Name(), "EventPrinter"}, new(RegisterReply))
	if err != nil {
		log.Fatal(err)

	}

	peer.Register(new(EventPrinter))
	peer.Serve()
}
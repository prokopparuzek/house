package main

import (
	"fmt"
	"log"

	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
)

type message struct {
	Timestamp   int64   `json:"timestamp"`
	Temperature float64 `json:"temperature"`
}

const subject = "room"

func main() {
	forever := make(chan bool)
	nc, err := nats.Connect("rpi3:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	sc, err := stan.Connect("measures", "reader1", stan.NatsConn(nc))
	if err != nil {
		panic(err)
	}
	defer sc.Close()
	log.Println("Connected")
	sc.Subscribe(subject, func(msg *stan.Msg) { fmt.Println(msg) }, stan.DurableName("bima"), stan.DeliverAllAvailable())
	<-forever
}

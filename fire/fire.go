package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	stan "github.com/nats-io/stan.go"
)

const subject = "room"

var scon stan.Conn
var client *firestore.Client
var ctx context.Context

func STANConnect(_ stan.Conn, _ error) {
	for true {
		time.Sleep(3 * time.Second)
		sc, err := stan.Connect("measures", "gun", stan.NatsURL("nats://rpi3:4222"), stan.SetConnectionLostHandler(STANConnect))
		if err == stan.ErrBadConnection {
			continue
		} else if err != nil {
			log.Panicln(err)
		} else {
			scon = sc
			break
		}
	}
}

func handleMsg(msg *stan.Msg) {
	var payload map[string]interface{}
	payload = make(map[string]interface{})
	ref := client.Collection("room-measures")
	json.Unmarshal(msg.Data, &payload)
	_, _, err := ref.Add(ctx, payload)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	forever := make(chan bool)
	STANConnect(nil, nil)
	defer scon.Close()
	log.Println("Connected")
	// firebase
	ctx = context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	scon.Subscribe(subject, handleMsg, stan.DurableName("1"), stan.DeliverAllAvailable())
	<-forever
}

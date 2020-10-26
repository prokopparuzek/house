package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	stan "github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
)

const subject = "room"
const logFile = "/var/log/fire.log"

var client *firestore.Client
var ctx context.Context

func handleMsg(msg *stan.Msg) {
	var payload map[string]interface{}
	payload = make(map[string]interface{})
	ref := client.Collection("room-measures")
	json.Unmarshal(msg.Data, &payload)
	log.WithField("payload", payload).Debug()
	_, _, err := ref.Add(ctx, payload)
	if err != nil {
		log.Println(err)
	}
	log.Debug("Fired")
}

func main() {
	log.SetOutput(os.Stderr)
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		log.WithField("file", logFile).Error(err)
	} else {
		log.SetOutput(f)
	}
	forever := make(chan bool)
	time.Sleep(3 * time.Second) // delay tpo start server
	sc, err := stan.Connect("measures", "gun", stan.NatsURL("nats://rpi3:4222"), stan.Pings(60, 1440))
	if err != nil {
		log.Panic(err)
	}
	defer sc.Close()
	log.Debug("Connected")
	// firebase
	ctx = context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	log.Debug("Connected to firebase")
	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	log.Debug("Connected to firestore")
	_, err = sc.Subscribe(subject, handleMsg, stan.DurableName("1"), stan.DeliverAllAvailable())
	if err != nil {
		log.Error(err)
	}
	log.Debug("Subscribed")
	<-forever
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	nats "github.com/nats-io/nats.go"
)

type data struct {
	Time        string  `json:"time"`
	Temperature float64 `json:"temperature"`
}

const Subject = "rpi2"
const MIME = "application/x-www-form-urlencoded"

func main() {
	// Připojení do gnatsd
	conn, err := nats.Connect("nats://rpi2:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("Connected")
	// Vytvoření kanálu
	econn, err := nats.NewEncodedConn(conn, nats.DEFAULT_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer econn.Close()
	channel := make(chan string)
	econn.BindRecvChan(Subject, channel)
	fmt.Println("Channel created!")
	for {
		var tmpStruct data
		var record [2]string
		var query string
		msg := <-channel
		json.Unmarshal([]byte(msg), &tmpStruct)
		record[0] = tmpStruct.Time
		record[1] = fmt.Sprintf("%f", tmpStruct.Temperature)
		query = "time=" + record[0] + "&temperature=" + record[1]
		// POST
		resp, err := http.Post(Url, MIME, bytes.NewBuffer([]byte(query)))
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		//fmt.Println(tmpStruct)

	}
}

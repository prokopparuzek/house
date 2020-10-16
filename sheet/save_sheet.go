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
	Temperature float64 `json:"temperature"`
}

type payload struct {
	Time string `json:"time"`
	Data data   `json:"data"`
}

const Subject = "rpi3"
const MIME = "application/x-www-form-urlencoded"

func main() {
	// Připojení do gnatsd
	conn, err := nats.Connect("nats://rpi3:4222")
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
		var tmpStruct payload
		var query string
		var temp_s string
		msg := <-channel
		json.Unmarshal([]byte(msg), &tmpStruct)
		temp_s = fmt.Sprintf("%f", tmpStruct.Data.Temperature)
		query = "time=" + tmpStruct.Time + "&temperature=" + temp_s
		// POST
		resp, err := http.Post(Url, MIME, bytes.NewBuffer([]byte(query)))
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		//fmt.Println(tmpStruct)

	}
}

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	nats "github.com/nats-io/nats.go"
)

type data struct {
	Time        string  `json:"time"`
	Temperature float64 `json:"temperature"`
}

const Subject = "rpi2"
const file_name = "temperature.csv"

func main() {
	// Otevření souboru s teplotou
	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		os.Create(file_name)
	}
	f, err := os.OpenFile(file_name, os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// Vytvoření cvs
	writer := csv.NewWriter(f)
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
		msg := <-channel
		json.Unmarshal([]byte(msg), &tmpStruct)
		record[0] = tmpStruct.Time
		record[1] = fmt.Sprintf("%f", tmpStruct.Temperature)
		writer.Write(record[:])
		writer.Flush()
		fmt.Println(tmpStruct)

	}
}

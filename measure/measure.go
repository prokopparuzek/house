package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	stan "github.com/nats-io/stan.go"
	cron "github.com/rk/go-cron"
)

type message struct {
	Timestamp   int64   `json:"timestamp"`
	Temperature float64 `json:"temperature"`
}

const device = "/sys/bus/w1/devices/28-03160521c6ff/temperature"
const csvFile = "/home/pi/data/measures.csv"
const subject = "room"

var scon stan.Conn

func getTemperature() float64 {
	file, err := os.Open(device)
	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()
	var temperature float64
	fmt.Fscan(file, &temperature)
	return temperature / 1000
}

func csvSave(msg *message) {
	f, err := os.OpenFile(csvFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	writer.Write([]string{fmt.Sprint(msg.Timestamp), fmt.Sprint(msg.Temperature)})
	writer.Flush()
	err = writer.Error()
	if err != nil {
		log.Panicln(err)
	}

}

func sendMeasures(_ time.Time) {
	var msg message
	var Jmsg []byte
	var err error
	var max int = 0
	msg.Temperature = getTemperature()
	msg.Timestamp = time.Now().Unix()
	//CSV
	csvSave(&msg)
	// NATS
	Jmsg, _ = json.Marshal(msg)
	for max = 0; max < 100; max++ {
		err = scon.Publish(subject, Jmsg)
		if err != nil {
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	if max >= 100 {
		log.Printf("Can't deliver %d\n", msg.Timestamp)
	}
}

func STANConnect(_ stan.Conn, _ error) {
	for true {
		time.Sleep(3 * time.Second)
		sc, err := stan.Connect("measures", "rpi3", stan.NatsURL("nats://rpi3:4222"), stan.SetConnectionLostHandler(STANConnect))
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

func main() {
	forever := make(chan bool)
	STANConnect(nil, nil)
	defer scon.Close()
	log.Println("Connected")
	// Cron
	cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, 00, 10, sendMeasures)
	cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, 05, 10, sendMeasures)
	cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, 10, 10, sendMeasures)
	cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, 15, 10, sendMeasures)
	cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, 20, 10, sendMeasures)
	cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, 25, 10, sendMeasures)
	cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, 30, 10, sendMeasures)
	cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, 35, 10, sendMeasures)
	cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, 40, 10, sendMeasures)
	cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, 45, 10, sendMeasures)
	cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, 50, 10, sendMeasures)
	cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, 55, 10, sendMeasures)
	<-forever
}

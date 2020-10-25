package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	stan "github.com/nats-io/stan.go"
	cron "github.com/rk/go-cron"
	log "github.com/sirupsen/logrus"
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
	log.WithField("file", device).Debug("Open file")
	var temperature float64
	count, err := fmt.Fscan(file, &temperature)
	log.Debug(temperature, count, err)
	return temperature / 1000
}

func csvSave(msg *message) {
	f, err := os.OpenFile(csvFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()
	log.WithField("file", csvFile).Debug("Open file")
	writer := csv.NewWriter(f)
	writer.Write([]string{fmt.Sprint(msg.Timestamp), fmt.Sprint(msg.Temperature)})
	writer.Flush()
	err = writer.Error()
	if err != nil {
		log.Panicln(err)
	}
	log.Debug("Stored csv")
}

func sendMeasures(_ time.Time) {
	var msg message
	var Jmsg []byte
	var err error
	var max int = 0
	var messageLoger *log.Entry = log.WithField("message", msg)
	msg.Temperature = getTemperature()
	msg.Timestamp = time.Now().Unix()
	messageLoger.Debug("Measure")
	//CSV
	csvSave(&msg)
	// NATS
	Jmsg, _ = json.Marshal(msg)
	for max = 0; max < 100; max++ {
		err = scon.Publish(subject, Jmsg)
		if err != nil {
			messageLoger.Trace("Error, will retry")
			time.Sleep(5 * time.Second)
		} else {
			messageLoger.Debug("Deliver")
			break
		}
	}
	if max >= 100 {
		messageLoger.Error("Cannot deliver")
	}
}

func STANConnect(_ stan.Conn, _ error) {
	for true {
		time.Sleep(3 * time.Second)
		sc, err := stan.Connect("measures", "rpi3", stan.NatsURL("nats://rpi3:4222"), stan.SetConnectionLostHandler(STANConnect))
		if err == stan.ErrBadConnection {
			log.Debug("Retrying")
			continue
		} else if err != nil {
			log.Panic(err)
		} else {
			scon = sc
			log.Debug("Connect")
			break
		}
	}
}

func main() {
	// logrus
	log.SetOutput(os.Stderr)
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	forever := make(chan bool)
	STANConnect(nil, nil)
	defer scon.Close()
	log.Debug("Connected-defer")
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
	log.Debug("Set CRON")
	<-forever
}

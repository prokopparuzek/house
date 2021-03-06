package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
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
const logFile = "/var/log/measure.log"
const subject = "room"

var scon stan.Conn

func getTemperature() float64 {
	var file *os.File
	var err error
	for {
		file, err = os.Open(device)
		if err != nil {
			log.Error(err)
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}
	defer file.Close()
	log.WithField("file", device).Debug("Open file")
	var temperature float64
	for true {
		count, err := fmt.Fscan(file, &temperature)
		if err == io.EOF {
			log.Debug("Retry read")
			continue
		} else if err != nil {
			log.Error(err)
			break
		} else {
			log.Debug(temperature, count, err)
			break
		}
	}
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
	var messageLoger *log.Entry = log.WithField("message", &msg)
	msg.Temperature = getTemperature()
	msg.Timestamp = time.Now().Unix()
	messageLoger.Debug("Measure")
	//CSV
	csvSave(&msg)
	// NATS
	Jmsg, _ = json.Marshal(msg)
	for max = 0; max < 120; max++ {
		err = scon.Publish(subject, Jmsg)
		if err != nil {
			messageLoger.Trace("Error, will retry")
			time.Sleep(60 * time.Second)
		} else {
			messageLoger.Debug("Deliver")
			break
		}
	}
	if max >= 120 {
		messageLoger.Error("Cannot deliver")
	}
}

func main() {
	// logrus
	log.SetOutput(os.Stderr)
	log.SetReportCaller(true)
	log.SetLevel(log.ErrorLevel)
	log.SetFormatter(&log.JSONFormatter{})
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		log.WithField("file", logFile).Error(err)
	} else {
		log.SetOutput(f)
	}
	forever := make(chan bool)
	var sc stan.Conn
	for {
		sc, err = stan.Connect("measures", "rpi3", stan.NatsURL("nats://rpi3:4222"), stan.Pings(60, 1440))
		if err != nil {
			log.Error(err)
			time.Sleep(time.Second * 30)
			continue
		}
		break
	}
	defer sc.Close()
	scon = sc
	log.Debug("Connected")
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
	//cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, cron.ANY, 10, sendMeasures)
	log.Debug("Set CRON")
	<-forever
}

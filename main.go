package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/splable/agent/v1/agent"
	"github.com/splable/agent/v1/conf"
)

const (
	// LightChannel reports sensor lux values.
	LightChannel = "LightChannel"
	// TemperatureChannel reports sensor temperature values.
	TemperatureChannel = "TemperatureChannel"
)

func main() {
	// Load configuration from splable.yml.
	var conf conf.File
	conf.GetConf()

	// Subscribe to web socket channels.
	agent := agent.NewSocket(conf)
	agent.Channel.Subscribe(TemperatureChannel)
	agent.Channel.Subscribe(LightChannel)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	for {
		select {
		case <-ticker.C:
			agent.Channel.ReportSHT3x(conf.Environment, TemperatureChannel)
			agent.Channel.ReportTSL2591(conf.Environment, LightChannel)
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupted, exiting...")
			ticker.Stop()
			agent.CloseSocket()
			return
		}
	}
}

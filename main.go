package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/splable/agent/v1/agent"
	"github.com/splable/agent/v1/conf"
	"github.com/splable/agent/v1/logger"
)

const (
	// HumidityChannel reports relative humidity values.
	HumidityChannel = "HumidityChannel"
	// LightChannel reports sensor lux values.
	LightChannel = "LightChannel"
	// TemperatureChannel reports sensor temperature values.
	TemperatureChannel = "TemperatureChannel"
	// PressureChannel reports sensor air pressure values.
	PressureChannel = "PressureChannel"
)

func main() {
	// Load configuration from splable.yml.
	var conf conf.File
	conf.GetConf()

	// Setup the logger.
	l := logger.CreateLogger(&conf)
	l.SetLevel(logger.NOTICE)

	welcomeMessage :=
		"\n" +
			"%s           _       _     _    \n" +
			" ___ _ __ | | __ _| |__ | | ___ \n" +
			"/ __| '_ \\| |/ _` | '_ \\| |/ _ \\\n" +
			"\\__ \\ |_) | | (_| | |_) | |  __/\n" +
			"|___/ .__/|_|\\__,_|_.__/|_|\\___|\n" +
			"    |_|\n%s\n"

	fmt.Fprintf(os.Stderr, welcomeMessage, "", "")
	l.Notice("Starting splable-agent with PID: %s", fmt.Sprintf("%d", os.Getpid()))

	// Subscribe to web socket channels.
	l.Info("Connecting to %s", conf.Hostname)
	agent := agent.NewSocket(l, conf)

	l.Info("Subscribing to humidity channel")
	agent.Channel.Subscribe(l, HumidityChannel)
	l.Info("Subscribing to temperature channel")
	agent.Channel.Subscribe(l, TemperatureChannel)
	l.Info("Subscribing to light channel")
	agent.Channel.Subscribe(l, LightChannel)
	l.Info("Subscribing to pressure channel")
	agent.Channel.Subscribe(l, PressureChannel)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	for {
		select {
		case <-ticker.C:
			l.Notice("Starting reporting loop...")
			agent.Channel.ReportSHT3x(l, conf, TemperatureChannel, HumidityChannel)
			agent.Channel.ReportTSL2591(l, conf, LightChannel)
			agent.Channel.ReportMPL3115A2(l, conf, PressureChannel)
		case <-done:
			return
		case <-interrupt:
			l.Notice("Interrupted, exiting...")
			ticker.Stop()
			agent.CloseSocket()
			return
		}
	}
}

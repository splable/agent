package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jimnelson2/tsl2591"
)

// MeasureTSL2591Rand gets the current visible light value from the TSL2591 sensor.
func MeasureTSL2591Rand() float64 {
	seed := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(seed)
	min := 0.0
	max := 50.0
	lux := min + rand.Float64()*(max-min)

	return lux
}

// MeasureTSL2591 gets the current visible light value from the TSL2591 sensor.
func MeasureTSL2591() float64 {
	tsl, err := tsl2591.NewTSL2591(&tsl2591.Opts{
		Gain:   tsl2591.GainMed,
		Timing: tsl2591.Integrationtime100MS,
	})
	if err != nil {
		log.Panic(err)
	}

	lux, err := tsl.Lux()
	if err != nil {
		log.Panic(err)
	}

	return lux
}

// ReportTSL2591 sends the current visible light value to the websocket channel.
func (c *ChannelService) ReportTSL2591(environment string, channelName string) {
	sensorValue := 0.0
	if environment == "development" {
		sensorValue = MeasureTSL2591Rand()
	} else {
		sensorValue = MeasureTSL2591()
	}

	content := channelContent{
		Datetime: time.Now().UTC().String(),
		Message:  fmt.Sprintf("%.2f", sensorValue),
	}

	identifier := channelIdentifier{
		Channel: channelName,
	}

	data := channelData{
		Action:  ActionName,
		Content: content,
	}

	encodedIdentifier, err := json.Marshal(identifier)
	if err != nil {
		log.Panic(err)
	}

	encodedData, err := json.Marshal(data)
	if err != nil {
		log.Panic(err)
	}

	message := ChannelMessage{
		Command:    CommandTypeMessage,
		Identifier: string(encodedIdentifier),
		Data:       string(encodedData),
	}

	encodedMessage, err := json.Marshal(message)
	if err != nil {
		log.Panic(err)
	}

	// TODO: Need to handle message failures.
	c.client.socket.SendText(string(encodedMessage))
}

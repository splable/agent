package agent

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/jimnelson2/tsl2591"
	"github.com/splable/agent/v1/conf"
	"github.com/splable/agent/v1/logger"
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
func MeasureTSL2591(l logger.Logger) float64 {
	tsl, err := tsl2591.NewTSL2591(&tsl2591.Opts{
		Gain:   tsl2591.GainMed,
		Timing: tsl2591.Integrationtime100MS,
	})
	if err != nil {
		l.Error("Error connecting to the TSL2591 sensor using the I2C bus: %s", err)
	}

	lux, err := tsl.Lux()
	if err != nil {
		l.Error("Error reading TSL2591 sensor value: %s", err)
	}

	return lux
}

// ReportTSL2591 sends the current visible light value to the websocket channel.
func (c *ChannelService) ReportTSL2591(l logger.Logger, conf conf.File, channelName string) {
	sensorValue := 0.0
	if conf.Environment == "development" {
		sensorValue = MeasureTSL2591Rand()
	} else {
		sensorValue = MeasureTSL2591(l)
	}

	l.Info("Light = %.2f", sensorValue)

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
		l.Error("Error sending sensor value to %s channel: %s", channelName, err)
	}

	encodedData, err := json.Marshal(data)
	if err != nil {
		l.Error("Error sending sensor value to %s channel: %s", channelName, err)
	}

	message := ChannelMessage{
		Command:    CommandTypeMessage,
		Identifier: string(encodedIdentifier),
		Data:       string(encodedData),
	}

	encodedMessage, err := json.Marshal(message)
	if err != nil {
		l.Error("Error sending sensor value to %s channel: %s", channelName, err)
	}

	// TODO: Need to handle message failures.
	c.client.socket.SendText(string(encodedMessage))
}

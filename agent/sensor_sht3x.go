package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	"github.com/d2r2/go-sht3x"
)

// MeasureSHT3xRand gets the current temperature value from the SHT3x sensor.
func MeasureSHT3xRand() float64 {
	seed := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(seed)
	min := 0.0
	max := 30.0
	temp := min + rand.Float64()*(max-min)

	return temp
}

// MeasureSHT3x gets the current temperature value from the SHT3x sensor.
func MeasureSHT3x() float64 {
	// Create new connection to i2c-bus on 1 line with address 0x44.
	// Use i2cdetect utility to find device address over the i2c-bus
	// ls /dev/i2c* to find out bus line.
	i2c, err := i2c.NewI2C(0x44, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer i2c.Close()

	logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
	logger.ChangePackageLogLevel("sht3x", logger.InfoLevel)

	sensor := sht3x.NewSHT3X()

	temp, _, err := sensor.ReadTemperatureAndRelativeHumidity(i2c, sht3x.RepeatabilityLow)
	if err != nil {
		log.Fatal(err)
	}

	return float64(temp)
}

// ReportSHT3x sends the current temperature value to the websocket channel.
func (c *ChannelService) ReportSHT3x(environment string, channelName string) {
	sensorValue := 0.0
	if environment == "development" {
		sensorValue = MeasureSHT3xRand()
	} else {
		sensorValue = MeasureSHT3x()
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

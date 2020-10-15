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
func MeasureSHT3xRand() (float64, float64) {
	tempSeed := rand.NewSource(time.Now().UnixNano())
	tempRand := rand.New(tempSeed)
	tempMin := 0.0
	tempMax := 30.0
	temp := tempMin + tempRand.Float64()*(tempMax-tempMin)

	rhSeed := rand.NewSource(time.Now().UnixNano())
	rhRand := rand.New(rhSeed)
	rhMin := 40.0
	rhMax := 80.0
	rh := rhMin + rhRand.Float64()*(rhMax-rhMin)

	return temp, rh
}

// MeasureSHT3x gets the current temperature and humidity values from the SHT3x sensor.
func MeasureSHT3x() (float64, float64) {
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

	temp, rh, err := sensor.ReadTemperatureAndRelativeHumidity(i2c, sht3x.RepeatabilityLow)
	if err != nil {
		log.Fatal(err)
	}

	return float64(temp), float64(rh)
}

// ReportSHT3x generates sensor values based on environment. In dev, we just use random numbers
// since access to actual sensors is not available.
func (c *ChannelService) ReportSHT3x(environment string, tempChannelName string, humidityChannelName string) {
	tempValue := 0.0
	humidityValue := 0.0
	if environment == "development" {
		tempValue, humidityValue = MeasureSHT3xRand()
	} else {
		tempValue, humidityValue = MeasureSHT3x()
	}

	ReportValue(c, tempChannelName, tempValue)
	ReportValue(c, humidityChannelName, humidityValue)
}

// ReportValue sends the current temperature and humidity values to websocket channels.
func ReportValue(c *ChannelService, channelName string, value float64) {
	content := channelContent{
		Datetime: time.Now().UTC().String(),
		Message:  fmt.Sprintf("%.2f", value),
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

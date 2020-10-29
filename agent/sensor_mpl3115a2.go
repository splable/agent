package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-mpl3115a2"
)

// MeasureMPL3115A2Rand gets the current visible light value from the TSL2591 sensor.
func MeasureMPL3115A2Rand() float64 {
	seed := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(seed)
	min := 10000.0
	max := 10200.0
	pressure := min + rand.Float64()*(max-min)

	return pressure
}

// MeasureMPL3115A2 gets the current visible light value from the TSL2591 sensor.
func MeasureMPL3115A2() float64 {
	// Create new connection to i2c-bus on 1 line with address 0x60.
	// Use i2cdetect utility to find device address over the i2c-bus
	i2c, err := i2c.NewI2C(0x60, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer i2c.Close()

	sensor := mpl3115a2.NewMPL3115A2()

	// Oversample Ratio - define precision, from low(0) to high(7)
	osr := 3
	pressure, _, err := sensor.MeasurePressure(i2c, osr)
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("Pressure = %v Pa, temperature = %v *C", p, t)

	// a, t2, err = sensor.MeasureAltitude(i2c, osr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("Altitude = %v m, temperature = %v *C", a, t2)

	return float64(pressure)
}

// ReportMPL3115A2 sends the current visible light value to the websocket channel.
func (c *ChannelService) ReportMPL3115A2(environment string, channelName string) {
	sensorValue := 0.0
	if environment == "development" {
		sensorValue = MeasureMPL3115A2Rand()
	} else {
		sensorValue = MeasureMPL3115A2()
	}

	content := channelContent{
		Datetime: time.Now().UTC().String(),
		Message:  fmt.Sprintf("%.0f", sensorValue),
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

package agent

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/d2r2/go-i2c"
	d2r2Logger "github.com/d2r2/go-logger"
	"github.com/d2r2/go-mpl3115a2"
	"github.com/splable/agent/v1/conf"
	"github.com/splable/agent/v1/logger"
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
func MeasureMPL3115A2(l logger.Logger) float64 {
	// Create new connection to i2c-bus on 1 line with address 0x60.
	// Use i2cdetect utility to find device address over the i2c-bus
	i2c, err := i2c.NewI2C(0x60, 1)
	if err != nil {
		l.Error("Error connecting to the MPL3115A2 sensor using the I2C bus: %s", err)
	}
	defer i2c.Close()

	d2r2Logger.ChangePackageLogLevel("i2c", d2r2Logger.InfoLevel)
	d2r2Logger.ChangePackageLogLevel("mpl3115a2", d2r2Logger.InfoLevel)

	sensor := mpl3115a2.NewMPL3115A2()

	// Oversample Ratio - define precision, from low(0) to high(7)
	osr := 3
	pressure, _, err := sensor.MeasurePressure(i2c, osr)
	if err != nil {
		l.Error("Error reading MPL3115A2 sensor value: %s", err)
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
func (c *ChannelService) ReportMPL3115A2(l logger.Logger, conf conf.File, channelName string) {
	sensorValue := 0.0
	if conf.Environment == "development" {
		sensorValue = MeasureMPL3115A2Rand()
	} else {
		sensorValue = MeasureMPL3115A2(l)
	}

	l.Info("Pressure = %.0f", sensorValue)

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

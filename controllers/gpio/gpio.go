package gpio

import (
	gbgpio "gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

var pinDrivers []*pinDriver

type pinDriver struct {
	Driver *gbgpio.DirectPinDriver
	Pin    string
}

func setupPins(conn *raspi.Adaptor, cfg *Config) {
	for pin := range cfg.Pins {
		pinDrivers = append(pinDrivers, &pinDriver{
			gbgpio.NewDirectPinDriver(conn, cfg.Pins[pin]),
			cfg.Pins[pin],
		})
	}
}

func allOn() {
	for p := range pinDrivers {
		pinDrivers[p].Driver.On()
	}
}

func allOff() {
	for p := range pinDrivers {
		pinDrivers[p].Driver.Off()
	}
}

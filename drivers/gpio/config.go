package gpio

import "gobot.io/x/gobot/platforms/raspi"

type Config struct {
	Conn *raspi.Adaptor
	Pins []string
}

package main

import (
	"github.com/thenrich/the-cabin-burned/controllers/gpio"
	"github.com/thenrich/the-cabin-burned/controllers/command"
	"gobot.io/x/gobot/platforms/raspi"
)

// Controller defines the behavior of a light controller
type Controller interface {
	// Start should begin a goroutine to monitor the state of the light
	Start()
	// Name should return the name of the light
	Name() string
	// State should return the state of the light
	State() int
	// Activate should activate the light
	Activate()
	// Deactivate should deactivate the light
	Deactivate()
}

// Lights is the parent structure that controls all light interactions
type Lights struct {
	config *LightsConfig
	lights map[string]Controller
}

// LightsConfig provides configuration for the main Lights structure
type LightsConfig struct {
	// Exclusive only allows one light in the config to be enabled
	// at a time
	Exclusive bool
}

// NewLights creates a new Lights object
func NewLights(c *LightsConfig) *Lights {
	l := &Lights{c, make(map[string]Controller)}
	return l
}

// AddLight adds a light controller
func (l *Lights) AddLight(c Controller) {
	l.lights[c.Name()] = c
}

// handleStateChange calls the appropriate activate/deactivate on each
// controller based on configuration settings
func (l *Lights) handleStateChange(light string, state int) {
	if l.config.Exclusive {
		if state == StateOn {
			for key := range l.lights {
				if key == light { // Skip if we're on this light
					continue
				}
				l.lights[key].Deactivate()
			}
			l.lights[light].Activate()
		} else if state == StateOff {
			l.lights[light].Deactivate()
		}
	} else {
		if state == StateOn {
			l.lights[light].Activate()
		} else if state == StateOff {
			l.lights[light].Deactivate()
		}
	}
}

func main() {
	l := NewLights(&LightsConfig{Exclusive: true})

	l.AddLight(NewControl("christmas_lights_music", command.NewLights("python", "/home/pi/lightshowpi/py/hardware_controller.py", "--state", "dance")))
	l.AddLight(NewControl("christmas_lights",
		gpio.NewLights(
			&gpio.Config{
				Conn: raspi.NewAdaptor(),
				Pins: []string{
					"11",
					"12",
					"13",
					"15",
					"16",
					"18",
					"7",
					"32",
				}})))

	m := NewMQTTHandler(l, &HandlerConfig{NamespacePrefix: "home/outside"})
	h := NewHTTPHandler(l, &HandlerConfig{NamespacePrefix: "/lights"})

	AddHandler(m)
	AddHandler(h)

	Listen()
}

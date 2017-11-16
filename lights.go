package main

import (
	"fmt"
	"github.com/thenrich/the-cabin-burned/drivers/gpio"
	"github.com/thenrich/the-cabin-burned/drivers/command"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"github.com/pkg/errors"
	"strings"
)

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
	config, err := ReadConfig("tcb.yaml")
	if err != nil {
		log.Fatal(errors.Wrap(err, "unable to read config, exiting..."))
	}

	if config.MQTT == nil {
		log.Fatal(errors.New("mqtt configuration is required, check config file"))
	}

	if config.Lights == nil {
		log.Fatal(errors.New("light configuration is required, check config file"))
	}

	l := NewLights(&LightsConfig{Exclusive: true})

	for lightOption := range config.Lights {
		if config.Lights[lightOption].Kind == "gpio" {
			if config.Lights[lightOption].Pins == "" {
				log.Fatal(errors.Errorf("invalid pin configuration for light: %s", config.Lights[lightOption].Name))
			}
			gpioCfg := &gpio.Config{
				Conn: raspi.NewAdaptor(),
				Pins: strings.Split(config.Lights[lightOption].Pins, ","),
			}
			l.AddLight(NewControl(
				config.Lights[lightOption].Name,
				gpio.NewLights(gpioCfg),
				config.MQTT))
		}

		if config.Lights[lightOption].Kind == "command" {
			if config.Lights[lightOption].Command == "" {
				log.Fatal(errors.Errorf("invalid command for light: %s", config.Lights[lightOption].Name))
			}
			fmt.Println(config.Lights[lightOption])
			l.AddLight(
				NewControl(
					config.Lights[lightOption].Name,
					command.NewLights(config.Lights[lightOption].Command,
						config.Lights[lightOption].CommandArgs...),
					config.MQTT))
		}
	}

	m := NewMQTTHandler(l,
		&MQTTHandlerConfig{HandlerConfig: HandlerConfig{
			NamespacePrefix: config.MQTT.Prefix,
		}, Broker: config.MQTT.Broker})
	AddHandler(m)

	if config.HTTP != nil {
		h := NewHTTPHandler(l,
			&HTTPHandlerConfig{HandlerConfig: HandlerConfig{
				NamespacePrefix: config.HTTP.Prefix,
			}, Listen: config.HTTP.Listen})
		AddHandler(h)
	}

	Listen()
}

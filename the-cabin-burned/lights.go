package the_cabin_burned

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	"github.com/thenrich/the-cabin-burned/drivers/command"
	"github.com/thenrich/the-cabin-burned/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"strings"
	"time"
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

func Start(config *Configuration) {
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
			ctl := NewControl(
				config.Lights[lightOption].Name,
				gpio.NewLights(gpioCfg))
			ctl.Start()
			l.AddLight(ctl)
		}

		if config.Lights[lightOption].Kind == "command" {
			fmt.Println("COMMAND")
			if config.Lights[lightOption].Command == "" {
				log.Fatal(errors.Errorf("invalid command for light: %s", config.Lights[lightOption].Name))
			}
			ctl := NewControl(
				config.Lights[lightOption].Name,
				command.NewLights(config.Lights[lightOption].Command,
					config.Lights[lightOption].CommandArgs...))
			ctl.Start()
			l.AddLight(ctl)
		}
	}

	if config.HTTP != nil {
		h := NewHTTPHandler(l,
			&HTTPHandlerConfig{HandlerConfig: HandlerConfig{
				NamespacePrefix: config.HTTP.Prefix,
			}, Listen: config.HTTP.Listen})

		AddHandler(h)
	}
	opts := mqtt.NewClientOptions()
	log.Printf("add %s\n", config.MQTT.Broker)
	opts.AddBroker(config.MQTT.Broker)
	mClient := mqtt.NewClient(opts)
	if tok := mClient.Connect(); tok.WaitTimeout(time.Second * 10) != true {
		log.Println("unable to connect to mqtt broker after 10s")
		log.Fatal(tok.Error())
	}

	m := NewMQTTHandler(l, &mClient, config.MQTT.Prefix)
	AddHandler(m)

	Listen()
	select {}
}

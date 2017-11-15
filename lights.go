package main

import (
	"net/http"
	"fmt"
	"log"
	"strings"
	"github.com/eclipse/paho.mqtt.golang"
)

type Controller interface {
	Start()
	Name() string
	State() int
	Activate()
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

func (l *Lights) handleHttpStateChange(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	light, state := parts[2], parts[3]
	newState := string2state[state]

	l.handleStateChange(light, newState)
}

func (l *Lights) handleMqttStateChange(client mqtt.Client, message mqtt.Message) {
	parts := strings.Split(message.Topic(), "/")
	light, state := parts[2], string(message.Payload())
	var newState int
	if state == "ON" {
		newState = StateOn
	} else {
		newState = StateOff
	}

	l.handleStateChange(light, newState)
}

func (l *Lights) ServeHTTP() {
	for key := range l.lights {
		http.HandleFunc(fmt.Sprintf("/lights/%s/on", key), l.handleHttpStateChange)
		http.HandleFunc(fmt.Sprintf("/lights/%s/off", key), l.handleHttpStateChange)
	}
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (l *Lights) SubscribeMQTT() {
	c := NewMQTTClient()
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for key := range l.lights {
			if token := c.Subscribe(fmt.Sprintf("home/outside/%s/set", key), 0, l.handleMqttStateChange); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
	}

}

func main() {
	l := NewLights(&LightsConfig{Exclusive: true})
	l.AddLight(NewControl("christmas_lights_music", NewCommandLights("sleep", "15")))
	l.AddLight(NewControl("christmas_lights", NewGPIOLights()))

	l.SubscribeMQTT()
	l.ServeHTTP()
}

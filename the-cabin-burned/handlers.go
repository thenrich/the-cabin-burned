package the_cabin_burned

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"strings"
	"net/http"
	"log"
)

var handlers []Handler

type HandlerConfig struct {
	// NamespacePrefix is the prefix before a handler's resource identifier
	// Ex. home/outside/x/y
	NamespacePrefix string
}

type MQTTHandlerConfig struct {
	HandlerConfig

	Broker string
}

type HTTPHandlerConfig struct {
	HandlerConfig

	Listen string
}

type Handler interface {
	Serve()
}

type HTTPHandler interface {
	Handler
	handleStateChange(http.ResponseWriter, *http.Request)
}

type MQTTHandler interface {
	Handler
	handleStateChange(client mqtt.Client, message mqtt.Message)
}

type MQTT struct {
	config *MQTTHandlerConfig
	lights *Lights
}

func (m *MQTT) Serve() {
	// Create mqtt client
	c := NewMQTTClient(
		NewMQTTClientConfig(m.config.Broker, m.config.NamespacePrefix))
		
	for key := range m.lights.lights {
		if token := c.Client.Subscribe(fmt.Sprintf("%s/%s/set", m.config.NamespacePrefix, key), 0, m.handleStateChange); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
}

func (m *MQTT) handleStateChange(client mqtt.Client, message mqtt.Message) {
	parts := strings.Split(message.Topic(), "/")
	light, state := parts[2], string(message.Payload())
	var newState int
	if state == "ON" {
		newState = StateOn
	} else {
		newState = StateOff
	}

	m.lights.handleStateChange(light, newState)
}

func NewMQTTHandler(lights *Lights, config *MQTTHandlerConfig) *MQTT {
	return &MQTT{config, lights}
}

type HTTP struct {
	config *HTTPHandlerConfig
	lights *Lights
}

func (h *HTTP) Serve() {
	for key := range h.lights.lights {
		http.HandleFunc(fmt.Sprintf("%s/%s/on", h.config.NamespacePrefix, key), h.handleStateChange)
		http.HandleFunc(fmt.Sprintf("%s/%s/off", h.config.NamespacePrefix, key), h.handleStateChange)
	}
	log.Fatal(http.ListenAndServe(h.config.Listen, nil))
}
func (h *HTTP) handleStateChange(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	light, state := parts[2], parts[3]
	newState := string2state[state]

	h.lights.handleStateChange(light, newState)
}
func NewHTTPHandler(lights *Lights, config *HTTPHandlerConfig) *HTTP {
	return &HTTP{config, lights}
}

// AddHandler adds a handler
func AddHandler(h Handler) {
	handlers = append(handlers, h)
}

// Listen starts serving all handlers
func Listen() {
	for h := range handlers {
		handlers[h].Serve()
	}
}
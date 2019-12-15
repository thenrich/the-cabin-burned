package the_cabin_burned

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strings"
)

type MQTT struct {
	prefix string
	client *mqtt.Client
	lights *Lights
}

func NewMQTTHandler(lights *Lights, client *mqtt.Client, prefix string) *MQTT {
	return &MQTT{prefix, client, lights}
}

func (m *MQTT) Serve() {
	client := *m.client
	for key := range m.lights.lights {
		topic := fmt.Sprintf("/%s/%s/%s", m.prefix, key, "command")
		log.Printf("subcribe to %s\n", topic)
		client.Subscribe(topic, 1, m.onMessage)
		client.Publish(fmt.Sprintf("/%s/%s/%s", m.prefix, key, "available"), 1, true, []byte("available"))

		go m.handleStateChanges(m.lights.lights[key])
	}
}

func (m *MQTT) handleStateChanges(c Controller) {
	client := *m.client
	for {
		select {
		case state := <- c.StateChannel():
			log.Printf("got state: %s for %s\n", state, c.Name())
			client.Publish(fmt.Sprintf("/%s/%s/%s", m.prefix, c.Name(), "state"), 1, true, []byte(state))
		}
	}
}

func (m *MQTT) onMessage(c mqtt.Client, msg mqtt.Message) {
	// Find light from topic
	parts := strings.Split(msg.Topic(),"/")
	light := parts[2]

	cmd := string(msg.Payload())
	log.Printf("set state of %s to %s\n", light, cmd)
	m.lights.handleStateChange(light, string2state[cmd])
	// Publish state change back to broker
	//(*m.client).Publish(fmt.Sprintf("/%s/%s/%s", m.prefix, light, "state"), 1, true, []byte(cmd))
	msg.Ack()
}
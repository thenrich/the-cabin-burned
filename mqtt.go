package main

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
)

func Publish(topic string, state string) {
	c := NewMQTTClient()
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic = fmt.Sprintf("home/outside/%s/state", topic)

	token := c.Publish(topic, 0, true, state)
	token.Wait()
}

func NewMQTTClient() mqtt.Client {
	options := mqtt.NewClientOptions()
	options.AddBroker("tcp://192.168.1.5:1883")
	return mqtt.NewClient(options)

}

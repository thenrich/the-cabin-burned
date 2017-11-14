package main

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
)

func Publish(topic string, state string) {
	options := mqtt.NewClientOptions()
	options.AddBroker("tcp://192.168.1.5:1883")
	c := mqtt.NewClient(options)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic = fmt.Sprintf("home/outside/%s/state", topic)

	token := c.Publish(topic, 0, true, state)
	token.Wait()
}

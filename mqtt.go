package main

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
)

func Publish(clientConfig *MQTTClientConfig, topic string, state string) {
	options := mqtt.NewClientOptions()
	options.AddBroker(clientConfig.Broker)
	c := mqtt.NewClient(options)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic = fmt.Sprintf("%s/%s/state", clientConfig.Prefix, topic)

	token := c.Publish(topic, 0, true, state)
	token.Wait()
}

type MQTTClientConfig struct {
	Broker string
	Prefix string
}

func NewMQTTClientConfig(broker string, prefix string) *MQTTClientConfig {
	return &MQTTClientConfig{
		Broker: broker,
		Prefix: prefix,
	}

}

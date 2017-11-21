package the_cabin_burned

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"time"
)

func Publish(c *MQTTClient, topic string, state string) {
	topic = fmt.Sprintf("%s/%s/state", c.Config.Prefix, topic)

	token := c.Client.Publish(topic, 0, true, state)
	token.WaitTimeout(time.Second * 5)
}

type MQTTClient struct {
	Client mqtt.Client
	Config *MQTTClientConfig
}

func NewMQTTClient(cfg *MQTTClientConfig) *MQTTClient {
	options := mqtt.NewClientOptions()
	options.AddBroker(cfg.Broker)
	options.AutoReconnect = true
	options.CleanSession = true
	options.SetMaxReconnectInterval(time.Second * 1)
	options.KeepAlive = time.Second * 5
	c := mqtt.NewClient(options)

	if token := c.Connect(); token.WaitTimeout(time.Second*10) && token.Error() != nil {
		panic(token.Error())
	}

	return &MQTTClient{c, cfg}
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

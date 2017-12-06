package the_cabin_burned

import (
	"github.com/eclipse/paho.mqtt.golang"
	"log"
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
	options.SetConnectionLostHandler(func(c mqtt.Client, err error){
		log.Println("------ MQTT Connection Lost")	
		log.Println(err.Error())
	})
	options.SetOnConnectHandler(func(c mqtt.Client){
		log.Println("------ Connected to MQTT")
		if cfg.ConnectedHandler != nil {
			fmt.Println("Calling ConnectedHandler")
			cfg.ConnectedHandler(c)
		}
	})
	c := mqtt.NewClient(options)

	if token := c.Connect(); token.WaitTimeout(time.Second*10) && token.Error() != nil {
		panic(token.Error())
	}

	return &MQTTClient{c, cfg}
}

type MQTTClientConfig struct {
	Broker string
	Prefix string

	ConnectedHandler func(client mqtt.Client)
}

func NewMQTTClientConfig(broker string, prefix string) *MQTTClientConfig {
	return &MQTTClientConfig{
		Broker: broker,
		Prefix: prefix,
	}

}

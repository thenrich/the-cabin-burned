package the_cabin_burned

import (
	"io/ioutil"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// MQTTConfig is the host to pub/sub for MQTT messages: tcp://192.168.1.5:1883
type MQTTConfig struct {
	// MQTT Broker in the form of tcp://192.168.1.5:1883
	Broker string
	// MQTT topic prefix
	Prefix string
}

// HTTPConfig is the host to listen on for HTTP commands: 0.0.0.0:8080
type HTTPConfig struct {
	// Listen on the host and port: 0.0.0.0:8080
	Listen string
	// HTTP path prefix
	Prefix string
}

// ConfigLightsOptions defines a light driver
type ConfigLightsOptions struct {
	// Name of the light
	Name string
	// Kind of driver: gpio, command
	Kind string
	// Pins to use for driving the light (when kind is "gpio")
	Pins string `yaml:"pins,omitempty"`
	// Command to run to drive the light (when kind is "command")
	Command string `yaml:"command,omitempty"`
	// Arguments to provide to Command (when kind is "command")
	CommandArgs []string `yaml:"command_args,omitempty"`
}

type Configuration struct {
	// Listeners defines the services that listen for commands: http, mqtt
	MQTT      *MQTTConfig           `yaml:"mqtt,omitempty"`
	HTTP      *HTTPConfig           `yaml:"http,omitempty"`
	Pins      string                `yaml:"pins,omitempty"`
	Lights    []ConfigLightsOptions `yaml:"lights,omitempty"`
}

func readConfig(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read config file %s", path)
	}

	return content, nil
}

func ReadConfig(path string) (*Configuration, error) {
	content, err := readConfig(path)
	if err != nil {
		return nil, errors.Wrap(err, "can't read config")
	}
	cfg := &Configuration{}
	yaml.Unmarshal(content, cfg)

	return cfg, nil

}

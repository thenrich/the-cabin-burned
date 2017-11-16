package cmd

import (
	"log"
	"gobot.io/x/gobot/platforms/raspi"
	"strings"
	"fmt"
	"github.com/thenrich/the-cabin-burned/drivers/command"
	"github.com/spf13/cobra"
	"github.com/pkg/errors"
	"github.com/thenrich/the-cabin-burned/drivers/gpio"
	"github.com/thenrich/the-cabin-burned/the-cabin-burned"
)

var cfgFile string

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

var RootCmd = &cobra.Command{
	Use: "the-cabin-burned",
	Run: func(cmd *cobra.Command, args[]string) {
		run(cfgFile)
	},
}

func run(configFile string) {
	config, err := the_cabin_burned.ReadConfig(configFile)
	if err != nil {
		log.Fatal(errors.Wrap(err, "unable to read config, exiting..."))
	}

	if config.MQTT == nil {
		log.Fatal(errors.New("mqtt configuration is required, check config file"))
	}

	if config.Lights == nil {
		log.Fatal(errors.New("light configuration is required, check config file"))
	}

	l := the_cabin_burned.NewLights(&the_cabin_burned.LightsConfig{Exclusive: true})

	for lightOption := range config.Lights {
		if config.Lights[lightOption].Kind == "gpio" {
			if config.Lights[lightOption].Pins == "" {
				log.Fatal(errors.Errorf("invalid pin configuration for light: %s", config.Lights[lightOption].Name))
			}
			gpioCfg := &gpio.Config{
				Conn: raspi.NewAdaptor(),
				Pins: strings.Split(config.Lights[lightOption].Pins, ","),
			}
			l.AddLight(the_cabin_burned.NewControl(
				config.Lights[lightOption].Name,
				gpio.NewLights(gpioCfg),
				config.MQTT))
		}

		if config.Lights[lightOption].Kind == "command" {
			if config.Lights[lightOption].Command == "" {
				log.Fatal(errors.Errorf("invalid command for light: %s", config.Lights[lightOption].Name))
			}
			fmt.Println(config.Lights[lightOption])
			l.AddLight(
				the_cabin_burned.NewControl(
					config.Lights[lightOption].Name,
					command.NewLights(config.Lights[lightOption].Command,
						config.Lights[lightOption].CommandArgs...),
					config.MQTT))
		}
	}

	m := the_cabin_burned.NewMQTTHandler(l,
		&the_cabin_burned.MQTTHandlerConfig{HandlerConfig: the_cabin_burned.HandlerConfig{
			NamespacePrefix: config.MQTT.Prefix,
		}, Broker: config.MQTT.Broker})
	the_cabin_burned.AddHandler(m)

	if config.HTTP != nil {
		h := the_cabin_burned.NewHTTPHandler(l,
			&the_cabin_burned.HTTPHandlerConfig{HandlerConfig: the_cabin_burned.HandlerConfig{
				NamespacePrefix: config.HTTP.Prefix,
			}, Listen: config.HTTP.Listen})
		the_cabin_burned.AddHandler(h)
	}

	the_cabin_burned.Listen()
}

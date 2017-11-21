package cmd

import (
	"log"
	"github.com/spf13/cobra"
	"github.com/pkg/errors"
	"github.com/thenrich/the-cabin-burned/the-cabin-burned"
	"fmt"
)

var cfgFile string

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

var RootCmd = &cobra.Command{
	Use: "the-cabin-burned",
	Run: func(cmd *cobra.Command, args []string) {
		run(cfgFile)
	},
}

func run(configFile string) {
	fmt.Println("Starting...")
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

	the_cabin_burned.Start(config)

}

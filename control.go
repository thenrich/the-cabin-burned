package main

import (
	"fmt"
	"log"
)

const (
	// State definitions
	StateOn  = 10
	StateOff = 11
)

// state2string converts states to strings
var state2string = map[int]string{
	StateOn:  "on",
	StateOff: "off",
}

var string2state = map[string]int {
	"on": StateOn,
	"off": StateOff,
}

type LightControl interface {
	On(chan bool)
	Off()
}

type Control struct {
	Driver LightControl

	name       string
	state      int
	activate   chan bool
	deactivate chan bool
}

func (c *Control) State() int {
	return c.state
}

func (c *Control) Activate() {
	fmt.Printf("Activating %s\n", c.name)
	select {
	case c.activate <- true:
		fmt.Println("Activate.....")
	default:
		fmt.Println("Last action pending...")
	}
}

func (c *Control) Deactivate() {
	fmt.Printf("Deactivating %s\n", c.name)
	select {
	case c.deactivate <- true:
		fmt.Println("Deactivate.....")
	default:
		fmt.Println("Last action pending...")
	}
}

func (c *Control) setState(state int) {
	c.state = state
}

func (c *Control) Name() string {
	return c.name
}

func (c *Control) Start() {
	// done channel lets drivers return a signal
	// that they've completed
	done := make(chan bool)

	for {
		select {
		case <-c.activate:
			if c.state == StateOn {
				log.Println("Already on")
				continue
			}
			fmt.Println("Enable")
			c.Driver.On(done)
			c.setState(StateOn)

		case <-c.deactivate:
			if c.state == StateOff {
				log.Println("Already off")
				continue
			}
			c.Driver.Off()
			c.setState(StateOff)

			fmt.Println("Disable")
		case <-done:
			log.Println("Command completed")
			c.setState(StateOff)

		}
	}
}

func NewControl(name string, driver LightControl) *Control {
	ctrl := &Control{
		Driver:     driver,
		name:       name,
		activate:   make(chan bool),
		deactivate: make(chan bool)}

	go ctrl.Start()
	return ctrl
}

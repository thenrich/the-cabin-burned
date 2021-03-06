package the_cabin_burned

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

var string2state = map[string]int{
	"on":  StateOn,
	"off": StateOff,
}

// Controller defines the behavior of a light controller
type Controller interface {
	// Start should begin a goroutine to monitor the state of the light
	Start()
	// Name should return the name of the light
	Name() string
	// State should return the state of the light
	State() int
	// Activate should activate the light
	Activate()
	// Deactivate should deactivate the light
	Deactivate()

	StateChannel() chan string
}

// LightControl defines the behavior of the driving of the actual light itself
type LightControl interface {
	// On should turn the light on and optionally send a message on the
	// chanel if the light turns off on its own
	On(chan bool)
	// Off should turn off the light
	Off()
}

// Control implements the Controller interface for driving a LightControl
type Control struct {
	Driver LightControl

	name         string
	state        int
	activate     chan bool
	deactivate   chan bool
	stateChannel chan string
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
	// Publish state

	//Publish(c.MQTTClient, c.name, state2string[c.state])
	c.stateChannel <- state2string[c.state]
}

func (c *Control) StateChannel() chan string {
	return c.stateChannel
}

func (c *Control) Name() string {
	return c.name
}

func (c *Control) Start() {
	go c.start()
}

func (c *Control) start() {
	log.Printf("Start %s control", c.name)
	// done channel lets drivers return a signal
	// that they've completed
	done := make(chan bool)

	for {
		log.Printf("Wait for communication on %s channel\n", c.name)
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
			fmt.Println("Disable")
			c.Driver.Off()
			c.setState(StateOff)

		case <-done:
			log.Println("Command completed")
			c.setState(StateOff)

		}
	}
}

func NewControl(name string, driver LightControl) *Control {
	ctrl := &Control{
		Driver:       driver,
		name:         name,
		state:        StateOff,
		activate:     make(chan bool),
		deactivate:   make(chan bool),
		stateChannel: make(chan string)}

	return ctrl
}

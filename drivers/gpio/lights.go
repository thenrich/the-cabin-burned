package gpio

import "fmt"

type Lights struct {
	cfg *Config
}

func (g *Lights) On(done chan bool) {
	fmt.Println("gpio On")
	setupPins(g.cfg.Conn, g.cfg)
	allOn()
}

func (g *Lights) Off() {
	fmt.Println("gpio Off")
	setupPins(g.cfg.Conn, g.cfg)
	allOff()
}

func NewLights(cfg *Config) *Lights {
	return &Lights{cfg}
}

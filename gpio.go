package main

import "fmt"

type GPIOLights struct {
}

func (g *GPIOLights) On(done chan bool) {
	fmt.Println("gpio On")

}

func (g *GPIOLights) Off() {
	fmt.Println("gpio Off")

}

func NewGPIOLights() *GPIOLights {
	return &GPIOLights{}
}
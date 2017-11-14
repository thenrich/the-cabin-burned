package main

import (
	"fmt"
	"os/exec"
	"context"
	"log"
)

// run the command to invoke the lights
func runCmd(ctx context.Context, done chan error) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "sleep", "15")
	go func() {
		cmd.Start()
		done <- cmd.Wait()
	}()
	return cmd
}

type CommandLights struct {
	cmd *exec.Cmd
}

func (c *CommandLights) On(complete chan bool) {
	done := make(chan error)
	go func() {
		fmt.Println("Block while waiting for on done channel")
		err := <- done
		if err != nil {
			log.Println(err.Error())
		}
		complete <- true
	}()
	ctx := context.Background()
	c.cmd = runCmd(ctx, done)
}

func (c *CommandLights) Off() {
	fmt.Println("try to kill process")
	if c.cmd != nil {
		c.cmd.Process.Kill()
	}
}

func NewCommandLights() *CommandLights {
	return &CommandLights{}
}
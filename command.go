package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"context"
	"log"
)

// run the command to invoke the lights
func runCmd(ctx context.Context, done chan error, command string, args... string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, command, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil
	}
	go func() {
		cmd.Start()
		out, _ := ioutil.ReadAll(stderr)
		fmt.Printf("%s\n", out)
		done <- cmd.Wait()
	}()
	return cmd
}

type CommandLights struct {
	Command string
	Args    []string
	cmd     *exec.Cmd
}

func (c *CommandLights) On(complete chan bool) {
	done := make(chan error)
	go func() {
		fmt.Println("Block while waiting for on done channel")
		err := <-done
		if err != nil {
			log.Println(err.Error())
		}
		complete <- true
	}()
	ctx := context.Background()
	c.cmd = runCmd(ctx, done, c.Command, c.Args...)
}

func (c *CommandLights) Off() {
	fmt.Println("try to kill process")
	if c.cmd != nil {
		c.cmd.Process.Kill()
	}
}

func NewCommandLights(command string, args... string) *CommandLights {
	return &CommandLights{Command: command, Args: args}
}

package command

import (
	"fmt"
	"io/ioutil"
	"syscall"
	"os/exec"
	"context"
	"log"
)

// run the command to invoke the lights
func runCmd(ctx context.Context, done chan error, command string, args... string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, command, args...)

	// Set process group so we can kill the process' children when we kill it
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil
	}
	go func() {
		cmd.Start()
		out, _ := ioutil.ReadAll(stderr)
		fmt.Printf("%s\n", out)
		out, _ = ioutil.ReadAll(stdout)
		fmt.Printf("%s\n", out)
		done <- cmd.Wait()
	}()
	return cmd
}

type Lights struct {
	Command string
	Args    []string
	cmd     *exec.Cmd
}

func (c *Lights) On(complete chan bool) {
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

func (c *Lights) Off() {
	fmt.Println("try to kill process")
	if c.cmd != nil {
		pgid, err := syscall.Getpgid(c.cmd.Process.Pid) 
		if err != nil {
			log.Printf("error getting process id: %s\n", err.Error())
			return
		}
		syscall.Kill(-pgid, 15)
	}
}

func NewLights(command string, args... string) *Lights {
	fmt.Println(command)
	fmt.Println(args)
	return &Lights{Command: command, Args: args}
}

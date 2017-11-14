package main

import (
	"net/http"
	"fmt"
	"log"
	"strings"
)

type Controller interface {
	Start()
	Name() string
	State() int
	Activate()
	Deactivate()
}

type Lights struct {
	config *LightsConfig
	lights map[string]Controller
}

type LightsConfig struct {
	// Exclusive only allows one light in the config to be enabled
	// at a time
	Exclusive bool
}

func NewLights(c *LightsConfig) *Lights {
	l := &Lights{c, make(map[string]Controller)}
	return l
}

func (l *Lights) AddLight(c Controller) {
	l.lights[c.Name()] = c
}

func (l *Lights) handleHttpStateChange(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	light, state := parts[2], parts[3]
	newState := string2state[state]

	if l.config.Exclusive {
		if newState == StateOn {
			for key := range l.lights {
				if key == light { // Skip if we're on this light
					continue
				}
				l.lights[key].Deactivate()
			}
			l.lights[light].Activate()
		} else if newState == StateOff {
			l.lights[light].Deactivate()
		}
	} else {
		if newState == StateOn {
			l.lights[light].Activate()
		} else if newState == StateOff {
			l.lights[light].Deactivate()
		}
	}

}

func (l *Lights) Serve() {
	for key := range l.lights {
		http.HandleFunc(fmt.Sprintf("/lights/%s/on", key), l.handleHttpStateChange)
		http.HandleFunc(fmt.Sprintf("/lights/%s/off", key), l.handleHttpStateChange)
	}
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	l := NewLights(&LightsConfig{Exclusive: true})
	l.AddLight(NewControl("syncro", NewCommandLights()))
	l.AddLight(NewControl("regular", NewGPIOLights()))

	l.Serve()
}

// set the state of name to state and publish the state change
//func setState(name string, state int) {
//	State[name] = state
//
//	// Publish
//	Publish(name, state2string[state])
//}

//func handleLights(enable chan bool, disable chan bool) {
//	ctx := context.Background()
//	var cmd *exec.Cmd
//
//	done := make(chan error)
//
//	for {
//		select {
//		case <-enable:
//			if State["eee"] == StateOn {
//				log.Println("Already on")
//				continue
//			}
//			fmt.Println("Enable")
//			setState("eee", StateOn)
//			cmd = runCmd(ctx, done)
//		case <-disable:
//			if State["eee"] == StateOff {
//				log.Println("Already off")
//				continue
//			}
//			if cmd != nil {
//				cmd.Process.Kill()
//			}
//			fmt.Println("Disable")
//			setState("eee", StateOff)
//		case <-done:
//			log.Println("Command completed")
//			setState("eee", StateOff)
//
//		}
//	}
//}

//func main() {
//	State = make(map[string]int)
//	State["eee"] = StateOff
//
//	enableLightsChannel := make(chan bool)
//	disableLightsChannel := make(chan bool)
//
//	go handleLights(enableLightsChannel, disableLightsChannel)
//
//	http.HandleFunc("/lights/on", func(w http.ResponseWriter, r *http.Request) {
//		select {
//		case enableLightsChannel <- true:
//			fmt.Fprintf(w, "Lights on\n")
//		default:
//			fmt.Fprintf(w, "Last action still pending\n")
//		}
//
//	})
//	http.HandleFunc("/lights/off", func(w http.ResponseWriter, r *http.Request) {
//		select {
//		case disableLightsChannel <- true:
//			fmt.Fprintf(w, "Lights off\n")
//		default:
//			fmt.Fprintf(w, "Last action still pending\n")
//		}
//
//	})
//	log.Fatal(http.ListenAndServe(":8080", nil))
//}

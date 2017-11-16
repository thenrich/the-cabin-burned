package the_cabin_burned

// Lights is the parent structure that controls all light interactions
type Lights struct {
	config *LightsConfig
	lights map[string]Controller
}

// LightsConfig provides configuration for the main Lights structure
type LightsConfig struct {
	// Exclusive only allows one light in the config to be enabled
	// at a time
	Exclusive bool
}

// NewLights creates a new Lights object
func NewLights(c *LightsConfig) *Lights {
	l := &Lights{c, make(map[string]Controller)}
	return l
}

// AddLight adds a light controller
func (l *Lights) AddLight(c Controller) {
	l.lights[c.Name()] = c
}

// handleStateChange calls the appropriate activate/deactivate on each
// controller based on configuration settings
func (l *Lights) handleStateChange(light string, state int) {
	if l.config.Exclusive {
		if state == StateOn {
			for key := range l.lights {
				if key == light { // Skip if we're on this light
					continue
				}
				l.lights[key].Deactivate()
			}
			l.lights[light].Activate()
		} else if state == StateOff {
			l.lights[light].Deactivate()
		}
	} else {
		if state == StateOn {
			l.lights[light].Activate()
		} else if state == StateOff {
			l.lights[light].Deactivate()
		}
	}
}

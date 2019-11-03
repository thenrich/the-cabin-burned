package the_cabin_burned

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var handlers []Handler

type JsonResponse struct {
	Active string `json:"active",omitempty`
}

type HandlerConfig struct {
	// NamespacePrefix is the prefix before a handler's resource identifier
	// Ex. home/outside/x/y
	NamespacePrefix string
}

type HTTPHandlerConfig struct {
	HandlerConfig
	Listen string
}

type Handler interface {
	Serve()
}

type HTTPHandler interface {
	Handler
	handleStateChange(http.ResponseWriter, *http.Request)
}


type HTTP struct {
	config *HTTPHandlerConfig
	lights *Lights
}

func (h *HTTP) Serve() {
	for key := range h.lights.lights {
		http.HandleFunc(fmt.Sprintf("/%s/%s", h.config.NamespacePrefix, key), h.handleRequest)
		//http.HandleFunc(fmt.Sprintf("/%s/%s", h.config.NamespacePrefix, key), h.handleRequest)
	}

	log.Fatal(http.ListenAndServe(h.config.Listen, nil))
}
func (h *HTTP) handleRequest(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	light := parts[2]

	if r.Method == "POST" {
		payload := &JsonResponse{}
		pb, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		if err := json.Unmarshal(pb, payload); err != nil {
			log.Println(err.Error())
			w.WriteHeader(500)
			return
		}

		log.Println(payload.Active)
		newState := string2state[payload.Active]

		h.lights.handleStateChange(light, newState)

		writeResponseStruct(w, &JsonResponse{Active: state2string[newState]})

		return
	}

	writeResponseStruct(w, &JsonResponse{Active: state2string[h.lights.lights[light].State()]})

	return

}

func writeResponseStruct(w http.ResponseWriter, i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	log.Println(string(b))
	w.Write(b)
}
func NewHTTPHandler(lights *Lights, config *HTTPHandlerConfig) *HTTP {
	return &HTTP{config, lights}
}

// AddHandler adds a handler
func AddHandler(h Handler) {
	handlers = append(handlers, h)
}

// Listen starts serving all handlers
func Listen() {
	for h := range handlers {
		go handlers[h].Serve()
	}
}
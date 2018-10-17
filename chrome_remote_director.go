package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"

	// for some reason gorilla websocket doesn't work with chrome
	// websocket when we use it as a client. we use go.net websockets
	// instead.
	//  "github.com/gorilla/websocket"
	"encoding/json"
	"io/ioutil"
)

type chromeResource struct {
	WebsocketDebuggerUrl string `json:"webSocketDebuggerUrl"`
	Type                 string `json:"type"`
	Url                  string `json:"url"`
	//all i care about for now.
}
type chromeCommand struct {
	ID     int               `json:"id"`
	Method string            `json:"method"`
	Params map[string]string `json:"params"`
}

type ChromeRemoteDirector struct {
	config         *CastConfig
	sessionUrl     string
	commandCounter int
}

func NewChromeRemoteDirector(config *CastConfig) *ChromeRemoteDirector {
	return &ChromeRemoteDirector{config: config}
}

func (p *ChromeRemoteDirector) Start() {
	chromeremote := fmt.Sprintf("%s/json", p.config.RemoteChrome)
	res, err := http.Get(chromeremote)
	if err != nil {
		log.Fatalf(`ChromeRemote: cannot find chrome remote api at %s: %s.
    Please launch chrome at the host machine this way:

      $ google-chrome --remote-debugging-port=9222\n\nThen provide chrome's

    Then specify your remote api address this way:

      $ <binary and existing flags> -chrome-remote=http://<host-machine-ip>:9222

    `, chromeremote, err)
	}
	text, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("ChromeRemote: %s", err)
	}

	resources := []chromeResource{}
	err = json.Unmarshal(text, &resources)

	if err != nil {
		log.Fatalf("ChromeRemote: can't understand chrome api: %s", err)
	}

	// find newtab
	socketurl := ""
	for _, resource := range resources {
		if resource.Type == "page" && resource.WebsocketDebuggerUrl != "" {
			socketurl = resource.WebsocketDebuggerUrl
			break
		}
	}
	if socketurl == "" {
		log.Fatalf("ChromeRemote: cannot find controllable page on your remote chrome")
	}
	log.Printf("ChromeRemote: session tab is at %s", socketurl)

	p.sessionUrl = socketurl
}

func (p *ChromeRemoteDirector) Stop() {
}

func (p *ChromeRemoteDirector) Open(tag string, url string) {
	p.runRemoteCommand(chromeCommand{Method: "Page.navigate", Params: map[string]string{"url": url}})
}

func (p *ChromeRemoteDirector) Close(tag string) {
}

func (p *ChromeRemoteDirector) runRemoteCommand(command chromeCommand) {
	command.ID = p.commandCounter + 1

	ws, err := websocket.Dial(p.sessionUrl, "", "http"+p.sessionUrl[2:])
	if err != nil {
		log.Fatalf("ChromeRemote: cannot communicate with chrome: %s\n%v", err)
	}

	defer ws.Close()

	cmd, err := json.Marshal(command)
	_, err = ws.Write(cmd)
	if err != nil {
		log.Fatalf("ChromeRemote: socket reponse %s")
	}
}

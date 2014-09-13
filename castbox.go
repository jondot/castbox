package main

import(
  "github.com/codegangsta/martini"
  "github.com/gorilla/websocket"
  "fmt"
  "net/http"
  "log"
  "io/ioutil"
  "encoding/json"
  "flag"
  "os/signal"
  "os"
)



var config *CastConfig
var appreg *AppRegistry
var director *AppDirector


//todo: replace with middleware?
func allowAllAccess(res http.ResponseWriter, req *http.Request)  {
    h := res.Header()
    h.Set("Access-Control-Allow-Origin", "*")
    h.Set("Access-Control-Allow-Method", "*")
    h.Set("Access-Control-Expose-Headers", "*")
}

func withApp(context martini.Context, params martini.Params, res http.ResponseWriter){
  app := appreg.Get(params["name"])
  if app == nil {
    http.Error(res, "No such application.", 404)
  }
  context.Map(app)
}

func wsHandshake(context martini.Context, w http.ResponseWriter, r *http.Request)  {
        if r.Method != "GET" {
                http.Error(w, "Method not allowed", 405)
                return
        }
        ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
        if _, ok := err.(websocket.HandshakeError); ok {
                http.Error(w, "Not a websocket handshake", 400)
                return
        } else if err != nil {
          http.Error(w, fmt.Sprintf("error: %s", err), 500)
                log.Println(err)
                return
        }
        context.Map(ws)
}






func deviceDescGet(res http.ResponseWriter, req *http.Request) string{
  h := res.Header()
  h.Set("Content-Type", "application/xml")
  h.Set("Application-URL", fmt.Sprintf("http://%s/apps", config.Host))
  buf := config.ExecuteTemplate("device_desc", config)
  return string(buf.Bytes())
}


func appsGet(res http.ResponseWriter, req *http.Request) {
  // TODO find *the first* running app and redirect to it.

  res.WriteHeader(204)
}

func appGet(res http.ResponseWriter, app *App) (int, string){
  res.Header().Set("Content-Type", "application/xml")
  director.RequestGet(app)

  buf := config.ExecuteTemplate("apps", app)
  return 200, string(buf.Bytes())
}
func appDelete(res http.ResponseWriter, req *http.Request, app *App) (int, string){
  director.RequestStop(app)

  res.Header().Set("Content-Type", "application/xml")
  buf := config.ExecuteTemplate("apps", app)
  return 200, string(buf.Bytes())
}

func appPost(res http.ResponseWriter, req *http.Request, app *App) (int, string){
  body, err := ioutil.ReadAll(req.Body)
  if err != nil {
    log.Printf("Error: cannot start app %s, cannot read body %s", err)
    return 404, ""
  }
  log.Printf("Client: cast! %s with data: %s", app.Name, body)

  director.Open(app, string(body))

  res.Header().Set("Location", fmt.Sprintf("http://%s/apps/%s/%s", config.Host, app.Name, app.RunningInstancePath))
  return 201, ""
}

func connectionPost(res http.ResponseWriter, app *App)(int, string){
  res.Header().Set("Content-Type", "application/json")

  buf := config.ExecuteTemplate("connection", map[string]interface{}{ "Config":config, "App":app })
  return 200, string(buf.Bytes())
  
}


// 2nd screen app will get this connection
func wsSession(ws *websocket.Conn, app *App) {
  log.Printf("Session: %v (Remote joining)", app.Name)
  app.RegisterRemote(ws)
}

// 1st screen app, after requesting 'connection' gets this resource.
func wsReceiver(ws *websocket.Conn, app *App) {
  log.Printf("Receiver: %v (Receiver joining)", app.Name)
  app.RegisterReceiver(ws)
}

// 1st screen loads and asks for connection.
func wsConnection(ws *websocket.Conn) {
  defer ws.Close()
  requestedApp := ""
  for {
    mtype, p, err := ws.ReadMessage()
    if err !=nil {
      log.Printf("Connection: %s", err)
      return
    }

    req := &ConnectionRequestMessage{}
    err = json.Unmarshal(p, req)
    if err !=nil {
      log.Printf("Connection: %s", string(p))
      continue
    }


    log.Printf("Connection: %s", req.Type)
    var msg interface{}
    switch req.Type {
      case "REGISTER":
        requestedApp = req.Name
        msg = NewChannelRequestMessage()
      case "CHANNELRESPONSE":
        msg = NewNewChannelMessage(fmt.Sprintf("ws://%s/receiver/%s", config.Host, requestedApp))
    }

    data, _ := json.Marshal(msg)
    err = ws.WriteMessage(mtype, data)

    if err != nil {
      log.Printf("Connection: %s", err)
      return
    }

  }

}


var castfile = flag.String("castfile", "Castfile", "A json file with app configuration")




func main() {
  flag.Parse()

  // set up machine and configuration
  machine := Machine{}
  addr, _ := machine.GetAddress()
  config = NewCastConfig(*castfile, addr)

  log.Printf("Starting on address: %s", config.Host)
  log.Printf("Device config: %v", config)

  // sync our apps
  appreg = NewAppRegistry(config)
  appreg.Sync()

  // start director and point to idle app
  director = &AppDirector{ config:config, apps: appreg }
  director.Start()

  // exit handlers and cleanups
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt) // CTRL-C
  go func(){
      for _ = range c {
        //in this case only interrupt
        director.Stop()
        os.Exit(0)
      }
  }()
  



  // set up Dial REST and websockets
  m := martini.Classic()

  //Dial REST
  m.Get("/ssdp/device-desc.xml", allowAllAccess,  deviceDescGet)
  m.Get("/apps",                 allowAllAccess, appsGet)
  m.Get("/apps/:name",           allowAllAccess, withApp, appGet)
  m.Post("/apps/:name",          allowAllAccess, withApp, appPost)
  m.Delete("/apps/:name/:instance",        allowAllAccess, withApp, appDelete)
  m.Post("/connection/:name",    allowAllAccess, withApp, connectionPost)

  //websockets
  m.Get("/session/:name",        wsHandshake, withApp, wsSession)
  m.Get("/receiver/:name",        wsHandshake, withApp, wsReceiver)
  m.Get("/connection",        wsHandshake, wsConnection)

  m.Get("/system/control", wsHandshake, func()(){

  })


  // set up discovery server
  discovery := NewDiscovery(config)
  go discovery.StartServer()

  // listen on the webs
  http.ListenAndServe(":8008", m)
}



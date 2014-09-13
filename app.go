package main

import (
  "log"
  "github.com/gorilla/websocket"
  "time"
  "sync/atomic"
)

var msgcnt int32 = 0

type App struct {
  Name string
  Url string
  UseChannel bool
  AllowEmptyPOSTData bool
  ConnectionSvcURL string
  RunningInstancePath string
  IsRunning bool
  RemotesHub *Hub
  ReceiversHub *Hub
  Channel chan *AppMessage
  LastActive time.Time
}

type AppMessage struct {
  Data []byte
  Mtype int
  Id int
  To *Hub
}

func (app *App) StartChannel(){
  appname := app.Name
  log.Printf("%s: Starting event channel", appname)
  defer log.Printf("%s: Event channel is closed.", appname)
  for {
    select {
      case msg, ok := <- app.Channel:
        if !ok {
          return
        }

        hub := msg.To
        from := hub.FromLabel
        to   := hub.ToLabel


        for app.IsRunning && ( len(app.ReceiversHub.Members) == 0  || len(app.RemotesHub.Members) == 0 ){
          log.Printf("%s: Waiting for members to join.", appname)
          time.Sleep(2000*time.Millisecond)
        }


        app.KeepAlive()
        for ws := range hub.Members {
          log.Printf("%s: SEND [%d] %s -> %s (of %d) [%s]", appname, msg.Id, from, to, len(hub.Members), string(msg.Data))
          err := ws.WriteMessage(msg.Mtype, msg.Data)
          if err != nil{
            log.Printf("%s: Error pumping message to %s: %s", app, to, err)
          }
        }
    }
  }
}

func (app *App) KeepAlive(){
  app.LastActive = time.Now()
}

func (app *App) StopChannel(){
  if app.Channel == nil {
    return
  }
  close(app.Channel)
  app.ReceiversHub.Close()
  app.RemotesHub.Close()
}

func (app *App) RequestStop(){
}

func (app *App) RequestInfo(){
  //for now just a liveness signal
  app.KeepAlive()
}

func (app *App) Start(){
  app.KeepAlive()

  if app.IsRunning {
    return
  }

  log.Printf("%s: Starting.", app.Name)
  app.IsRunning = true
  app.Channel = make(chan *AppMessage, 1000)
  app.RemotesHub =   &Hub{ FromLabel: "receiver", ToLabel: "remote", Members: map[*websocket.Conn]bool{} }
  app.ReceiversHub = &Hub{ FromLabel: "remote", ToLabel: "receiver", Members: map[*websocket.Conn]bool{}}
  go app.StartChannel()
}

func (app *App) Stop(){
  if app.IsRunning == false {
    return
  }

  log.Printf("%s: Closing.", app.Name)
  app.IsRunning = false
  app.StopChannel()
}

func (app *App) removeSocket(ws *websocket.Conn, socketmap map[*websocket.Conn]bool){
    delete(socketmap, ws)
    ws.Close()
}

func (app *App) registerSocket(ws *websocket.Conn, hub *Hub, toHub *Hub){
  defer hub.RemoveMember(ws)
  hub.AddMember(ws)

  app.KeepAlive()

  from := toHub.FromLabel
  to   := toHub.ToLabel
  appname := app.Name

  for {
    mtype, p, err := ws.ReadMessage()
    if err !=nil {
      log.Printf("%s: Connection (%s): %s", app.Name, from, err)
      break
    }

    msg := &AppMessage{ Data: p,
                        Mtype: mtype,
                        Id: int(atomic.AddInt32(&msgcnt, 1)),
                        To: toHub }
    log.Printf("%s: QUEUE [%d] %s -> %s (of %d) [%s]", appname, msg.Id, from, to, len(hub.Members), string(msg.Data))
    app.Channel <- msg
  }
}

func (app *App) RegisterRemote(ws *websocket.Conn){
  // remotes -> rec
  app.registerSocket(ws, app.RemotesHub, app.ReceiversHub)
}

func (app *App) RegisterReceiver(ws *websocket.Conn){
  // rec -> remotes
  app.registerSocket(ws, app.ReceiversHub, app.RemotesHub)
}









type Hub struct {
  Members map[*websocket.Conn]bool
  FromLabel string
  ToLabel string
  App *App
}

func (hub *Hub) Close(){
  for ws := range hub.Members {
    hub.RemoveMember(ws)
  }
}

func (hub *Hub) RemoveMember(ws *websocket.Conn){
    delete(hub.Members, ws)
    ws.Close()
}

func (hub *Hub) AddMember(ws *websocket.Conn){
  hub.Members[ws] = true
}



package main

import(
  "strings"
  "time"
  "log"
)


type AppDirector struct {
  config *CastConfig
  webdirector Director
  apps *AppRegistry
}

func (dir *AppDirector) Start(){
  dir.webdirector = dir.decideAndCreateDirector()

  dir.webdirector.Start()
  idleapp := dir.apps.GetIdleApp()
  dir.Open(idleapp, "")

  go dir.ScavengeDeadApps()
}

func (dir *AppDirector) Stop(){
  dir.webdirector.Stop()
}

func (dir *AppDirector) RequestGet(app *App){
  app.RequestInfo()
}


func (dir *AppDirector) RequestStop(app *App){
  app.RequestStop()
}

func (dir *AppDirector) Open(app *App, body string){
  appurl := strings.Replace(app.Url, "${POST_DATA}", body, -1)
  app.Start()
  dir.webdirector.Open(app.Name, appurl)

  // stop all other apps
  for _,otherapp := range dir.apps.Applications {
    if otherapp == app {
      continue
    }
    otherapp.Stop()
  }
}


//
// Scavenger
//
// The scavenger should look at all the apps and decide if we're
// inactive.
//
// Inactivity is vaguely defined by the chromecast guideline as 'no
// receiver or remote is connected' and/or 'no activity in the app'.
//
// Once that happens it should kill/stop all inactive apps and go back
// to idleapp.
//
// Since we can't do vague in code, here's what we'll actually do :)
//
// an app is alive if:
//
// - it was just started
// - someone is continuously requesting info via GET on REST Dial
// - someone just joined on websockets (either side)
// - both sides is doing websockets communication (remotes, receivers)
//
// here's how it applies to real apps:
//
// 1. Youtube-type receiver
// - we know that there's no communication via sockets so we can't keepalive through that.
// - but we're able to signal life via REST activity:
//      * GET apps/YouTube
// - here, DELETE might be actually important so we need to take
// 'RequestStop' into account.
//
// 2. Regular receiver
// - sockets should have session and receivers ping-ponging. Each of those is
//   a liveness signal.
//
func (dir *AppDirector) ScavengeDeadApps(){
  idleapp := dir.apps.GetIdleApp()

  for {
    log.Printf("AppDirector: Scavenging dead apps")
    allInactive := true
    for _,app := range dir.apps.Applications {
      if app.IsRunning == false || app == idleapp {
        continue
      }
      activityDuration := time.Now().Sub(app.LastActive)
      if activityDuration < dir.config.IdleTime{
        allInactive = false
        log.Printf("  - ACTIVE %s\t\t\t%v", app.Name, activityDuration)
      }else{
        log.Printf("  - DEAD   %s\t\t\t%v", app.Name, activityDuration)
      }

    }

    if allInactive && !idleapp.IsRunning {
      log.Printf("AppDirector: all apps are dead. Going back to idle app.")
      dir.Open(idleapp, "")
    }

    time.Sleep(10*time.Second)
  }
}


func (dir *AppDirector) decideAndCreateDirector() Director {
  var director Director

  if dir.config.RemoteChrome != "" {
    director = NewChromeRemoteDirector(dir.config)
  }else{
    director = NewChromeDriverDirector(dir.config)
  }
  return director
}


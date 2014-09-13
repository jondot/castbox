package main
import (
  "net/http"
  "strings"
  "io/ioutil"
  "encoding/json"
  "fmt"
  "log"
)

type ChromeApp struct {
  Name string                `json:"app_name"`
  Url string                 `json:"url"`
  UseChannel bool            `json:"use_channel"`
  AllowEmptyPOSTData bool    `json:"allow_empty_post_data"`
}


type ChromeConfig struct {
  Applications []*ChromeApp               `json:"applications"`
  Configuration map[string]interface{}    `json:"configuration"`
}


type AppRegistry struct {
  Applications map[string]*App
  IdleAppID string
  config *CastConfig
}

func NewAppRegistry(config *CastConfig) (*AppRegistry){
  reg := &AppRegistry{ Applications: map[string]*App{}, config: config }
  return reg
}


func (reg *AppRegistry) Sync() {
  reg.injestApps(reg.fetchGoogleApps())
  reg.injestApps(reg.fetchCastfileApps())
  log.Printf("Apps: Loaded the following apps:")
  for _,app := range reg.Applications {
    log.Printf("  - %s (%s)", app.Name, app.Url)
  }
  log.Printf("  Idle app: %s", reg.GetIdleApp().Name)
}

func (reg *AppRegistry) fetchGoogleApps() string{
  log.Printf("Apps: sync'ing Google apps...")
  res, err := http.Get("https://clients3.google.com/cast/chromecast/device/config")
  if err != nil {
    log.Printf("Apps: cannot fetch from google: %s", err)
    if len(reg.Applications) == 0 {
      log.Fatal("Apps: Exiting")
    }
	}
	appstext, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
  return strings.Replace(string(appstext), ")]}'", "", -1) // WAT?
}

func (reg *AppRegistry) fetchCastfileApps() string{
  log.Printf("Apps: sync'ing Castfile apps...")
  content, err := ioutil.ReadFile(reg.config.Castfile)
  if err != nil{
    return ""
  }
  return string(content)
}

func (reg *AppRegistry) injestApps(appsjson string){
  if appsjson == ""{
    return
  }
  conf := &ChromeConfig{}
  err := json.Unmarshal([]byte(appsjson), conf)
  if err != nil {
    log.Fatalf("Apps: cannot parse google's json: %s", err)
  }

  for _,app := range conf.Applications {
    if app.Url != "" {
      reg.Applications[app.Name] = &App{
        Name: app.Name,
        Url: app.Url,
        UseChannel: app.UseChannel,
        RunningInstancePath: "run",
        AllowEmptyPOSTData: app.AllowEmptyPOSTData,
        ConnectionSvcURL: fmt.Sprintf("http://%s/connection/%s", reg.config.Host, app.Name),
      }
    }
  }
  reg.IdleAppID = conf.Configuration["idle_screen_app"].(string)
}



func (reg *AppRegistry) Get(name string) *App{
  return reg.Applications[name]
}

func (reg *AppRegistry) GetIdleApp() *App{
  return reg.Applications[reg.IdleAppID]
}




package main

import (
  "text/template"
  "bytes"
  "time"
  "encoding/json"
  "log"
  "io/ioutil"
  "fmt"
)

type CastFileDeviceConfig struct {
  Device struct {
    UUID string   `json:"uuid"`
    IdleMins int     `json:"idle_time_min"`
    FriendlyName string  `json:"name"`
    RemoteChrome string  `json:"remote_chrome"`
    ForceHost string  `json:"force_host"`
    ForceChromebin string  `json:"force_chromebin"`
  } `json:"device"`
}

type CastConfig struct {
  Templates map[string]*template.Template
  UUID string
  IP string
  Host string
  FriendlyName string
  Castfile string
  IdleTime time.Duration
  RemoteChrome string
  ForceChromebin string
}


func NewCastConfig(castfile string, host string) (*CastConfig){
  content, err := ioutil.ReadFile(castfile)
  if err != nil {
    log.Fatalf("Cannot find Castfile: %s", err)
  }

  cfdevice := &CastFileDeviceConfig{}
  err = json.Unmarshal([]byte(content), cfdevice)
  if err != nil {
    log.Fatalf("Cannot parse Castfile: %s", err)
  }
  device := cfdevice.Device
  

  templates := map[string]*template.Template{}

  templates["apps"] = template.Must(template.New("apps").Parse(TMPL_APPS))
  templates["connection"] = template.Must(template.New("connection").Parse(TMPL_CONNECTION))
  templates["device_desc"] = template.Must(template.New("device_desc").Parse(TMPL_DEVICE_DESC))
  templates["dial_response"] = template.Must(template.New("dial_response").Parse(TMPL_DIAL_RESPONSE))

  if device.ForceHost != "" {
    host = device.ForceHost
  }

  host = fmt.Sprintf("%s:8008", host)

  return &CastConfig{ UUID: device.UUID,
                      FriendlyName: device.FriendlyName,
                      IdleTime: time.Duration(device.IdleMins)*time.Minute,
                      RemoteChrome: device.RemoteChrome,
                      Templates:templates,
                      ForceChromebin: device.ForceChromebin,
                      IP: host,
                      Host: host,
                      Castfile: castfile  }
}

func (cc *CastConfig) ExecuteTemplate(name string, params interface{}) *bytes.Buffer{
  buf := &bytes.Buffer{}
  cc.Templates[name].Execute(buf, params)

  return buf
}

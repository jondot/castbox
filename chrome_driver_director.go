package main

import (
  "log"
  "os"
  "runtime"
  "io/ioutil"
  "fmt"
  "time"
  "net"
)

type ChromeDriverDirector struct {
  process *os.Process
  remoteDirector *ChromeRemoteDirector
  config *CastConfig
}


func NewChromeDriverDirector(config *CastConfig)(*ChromeDriverDirector){
  return &ChromeDriverDirector{ remoteDirector: NewChromeRemoteDirector(config), config: config }
}

func (p *ChromeDriverDirector) Start(){
  selectedPath := ""

  if config.ForceChromebin != "" {
    selectedPath = config.ForceChromebin
  } else {
    chromepaths := getChromePaths()
    if len(chromepaths) == 0 {
      log.Fatalf("ChromeDriver: not familiar with chrome on this OS")
    }

    for _, path := range(chromepaths){
      if fileExists(path){
        selectedPath = path
        break
      }
    }
  }

  if selectedPath == "" {
    log.Fatalf("ChromeDriver: could not locate chrome on your machine")
  }

  tempuserdir, err := ioutil.TempDir("","castbox")
  if err !=nil{
    log.Fatalf("ChromeDriver: cannot allocate temp dir - permissions or disk full?")
  }

  remotePort := "9515"
  chromeArgs := []string{
    "--userdir=" + tempuserdir,

    "--remote-debugging-port=" + remotePort,
    // debug mode, remote url, incognito, no startup etc.
  }

  devNull, err := os.Open(os.DevNull) // For read access.
  if err != nil {
    log.Fatal(err)
  }
  attrs := os.ProcAttr{Files: []*os.File{os.Stdin, devNull, devNull}}
  proc, err := os.StartProcess(selectedPath, chromeArgs, &attrs)
  if err != nil{
    log.Fatalf("ChromeDriver: could not start chrome in remote mode")
  }

  p.process = proc

  
  log.Printf("ChromeDriver: remote chrome is listening on port %s", remotePort)
  p.config.RemoteChrome = fmt.Sprintf("http://localhost:%s",remotePort)

  for {

    _, err := net.Dial("tcp", "localhost:"+remotePort)
    if err == nil {
      break
    }
    log.Printf("ChromeDriver: polling for chrome - not found yet. Kill me if this never ends.")
    time.Sleep(1*time.Second) //yuck yuck yuck
  }
  p.remoteDirector.Start()
}

func (p *ChromeDriverDirector) Stop(){
  p.process.Kill()
}

func (p *ChromeDriverDirector) Open(tag string, url string) {
  p.remoteDirector.Open(tag, url)
}

func (p *ChromeDriverDirector) Close(tag string){
}

//c:\Users\dotan\AppData\Local\Google\Chrome\Application\chrome.exe

func getChromePaths() []string {
  homedir := os.Getenv("HOME")
  switch runtime.GOOS {
    case "linux":
      return []string{
        "/usr/bin/google-chrome",
        "/opt/google/chrome/google-chrome",
        "/usr/local/bin/google-chrome",
        "/usr/local/sbin/google-chrome",
        "/bin/google-chrome",
        "/sbin/google-chrome",
      }
    case "windows":
      return []string{
        homedir + `\AppData\Local\Google\Chrome\Application\chrome.exe`,
	`c:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
      }
    case "darwin":
      return []string{
        "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
      }
    default:
      return []string{ }
  }
}

func fileExists(name string) bool {
  if _, err := os.Stat(name); err != nil {
    if os.IsNotExist(err) {
      return false
    }
  }
  return true
}




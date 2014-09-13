package main

import (
  "net"
  "log"
  "regexp"
)


type Discovery struct {
  config *CastConfig
}

type DialResponse struct {
  IP string
  UUID string
}


func NewDiscovery(config *CastConfig) (*Discovery){
  return &Discovery{ config:config }
}

func (d *Discovery) StartServer(){
  //read dial_response -> remember line ends need to be \r\n
  //TODO run template with params:
  // server's ip -> might be tricky. might need to establish connection to requesting client and then extract our own ip.
  // uuid

  msearch := regexp.MustCompile(`(?s)M-SEARCH.*urn:dial-multiscreen-org:service:dial:1.*`)
  udpAddr, err := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
  if err != nil {
          log.Fatalf("ResolveUDPAddr failed: %s\n", err)
  }
  socket, err := net.ListenMulticastUDP("udp", nil, udpAddr)
  if err != nil {
          log.Fatalf("ListenUDP failed: %s\n", err.Error())
  }

  log.Printf("Discovery: server started.")
  for {
          message := make([]byte, 4096)
          n, caddr, err := socket.ReadFromUDP(message)
          if err !=nil {
            log.Printf("Discovery: ERROR reading - %s", err)
            continue
          }

          if !msearch.MatchString(string(message)){
            continue
          }

          log.Printf("Discovery: client(%s) M-SEARCH %d bytes.", caddr, n)
          //log.Printf("\n--->\n%v---<\n", string(message))
          /*
            M-SEARCH * HTTP/1.1
            HOST: 239.255.255.250:1900
            MAN: "ssdp:discover"
            MX: 1
            ST: urn:dial-multiscreen-org:service:dial:1
          */

          // build response
          buf := d.config.ExecuteTemplate("dial_response", DialResponse{IP: d.config.Host, UUID: d.config.UUID})

          n, err = socket.WriteToUDP(buf.Bytes(), caddr)
          if err !=nil {
            log.Printf("Discovery: ERROR sending - %s", err)
            continue
          }
          //log.Printf("Discovery: client(%s) sent %d bytes", caddr, n)
          //log.Printf("\n--->\n%v---<\n", string(buf.Bytes()))
  }

  log.Printf("Discovery: server terminated.")
}


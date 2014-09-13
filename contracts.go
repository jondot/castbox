package main

type Director interface {
  Start()
  Open(string, string)
  Stop()
}

type ConnectionRequestMessage struct {
  Type string                `json:"type"`
  Name string                 `json:"name"`
}

type ChannelRequestMessage struct {
  Type string                `json:"type"`
  Sender int                `json:"senderId"`
  Request int               `json:"requestId"`
}
func NewChannelRequestMessage()(*ChannelRequestMessage){
  return &ChannelRequestMessage{ Type: "CHANNELREQUEST", Sender:1, Request: 1 }
}

type NewChannelMessage struct {
  Type string                `json:"type"`
  Sender int                `json:"senderId"`
  Request int               `json:"requestId"`
  URL string                `json:"URL"`
}
func NewNewChannelMessage(url string)(*NewChannelMessage){
  return &NewChannelMessage{ Type: "NEWCHANNEL", URL: url }
}




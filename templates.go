package main




var TMPL_APPS = `<?xml version='1.0' encoding='UTF-8'?>
<service xmlns='urn:dial-multiscreen-org:schemas:dial'>
    <name>{{ .Name }}</name>
    <options allowStop='true'/>

    {{ if .IsRunning }}
      <servicedata xmlns='urn:chrome.google.com:cast'>
          <connectionSvcURL>{{ .ConnectionSvcURL }}</connectionSvcURL>
          <protocols>
              <protocol>ramp</protocol>
          </protocols>
      </servicedata>
      <activity-status xmlns="urn:chrome.google.com:cast">
        <description>{{ .Name }} Receiver</description>
      </activity-status>
      <link rel='run' href='{{ .RunningInstancePath }}'/>
    {{ end }}

    {{ if .IsRunning }}
      <state>running</state>
    {{ else }}
      <state>stopped</state>
    {{ end }}
</service>
`

var TMPL_CONNECTION = `{
  "URL":"ws://{{.Config.Host}}/session/{{.App.Name}}?1",
  "pingInterval": 3
}
`

var TMPL_DEVICE_DESC = `<?xml version="1.0" encoding="utf-8"?>
<root xmlns="urn:schemas-upnp-org:device-1-0" xmlns:r="urn:restful-tv-org:schemas:upnp-dd">
    <specVersion>
      <major>1</major>
      <minor>0</minor>
    </specVersion>
    <URLBase>http://{{.Host}}</URLBase>
    <device>
        <deviceType>urn:schemas-upnp-org:device:dial:1</deviceType>
        <friendlyName>{{ .FriendlyName }}</friendlyName>
        <manufacturer>Google Inc.</manufacturer>
        <modelName>Eureka Dongle</modelName>
        <UDN>uuid:{{ .UUID }}</UDN>
        <serviceList>
            <service>
                <serviceType>urn:schemas-upnp-org:service:dial:1</serviceType>
                <serviceId>urn:upnp-org:serviceId:dial</serviceId>
                <controlURL>/ssdp/notfound</controlURL>
                <eventSubURL>/ssdp/notfound</eventSubURL>
                <SCPDURL>/ssdp/notfound</SCPDURL>
            </service>
        </serviceList>
    </device>
</root>
`


// note - HTTP spec, \r\n ends line and last space is meaningful.
var TMPL_DIAL_RESPONSE ="HTTP/1.1 200 OK\r\n"+
                        "LOCATION: http://{{.IP}}/ssdp/device-desc.xml\r\n"+
                        "CACHE-CONTROL: max-age=1800\r\n"+
                        "CONFIGID.UPNP.ORG: 7337\r\n"+
                        "BOOTID.UPNP.ORG: 7337\r\n"+
                        "USN: uuid:{{.UUID}}\r\n"+
                        "ST: urn:dial-multiscreen-org:service:dial:1\r\n"+
                        "\r\n"



package main

import (
  "net"
  "os"
  "errors"
)


type Machine struct {
}

func (m *Machine) GetAddress() (string, error){
    name, err := os.Hostname()
    if err != nil {
        return "", err
    }

    addrs, err := net.LookupHost(name)

    if err != nil {
        return "", err
    }

    for _, a := range addrs {
      return a, nil
    }

    return "", errors.New("cannot find any address")
}



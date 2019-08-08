package main

import (
	"crypto/rand"
	"encoding/hex"
	consul "github.com/hashicorp/consul/api"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

const (
	instanceIDFile      = "reko-id"
	instanceServiceName = "reko"
)

var (
	instanceID        string
	instanceServiceID string
	instanceCheckID   string
)

func ensureInstanceID() (err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(instanceIDFile); err != nil {
		if !os.IsNotExist(err) {
			return
		}
	}
	if len(buf) == 0 {
		id := make([]byte, 12, 12)
		if _, err = rand.Read(id); err != nil {
			return
		}
		instanceID = hex.EncodeToString(id)
		if err = ioutil.WriteFile(instanceIDFile, []byte(instanceID), 0644); err != nil {
			return
		}
	} else {
		instanceID = strings.TrimSpace(string(buf))
	}

	instanceServiceID = "reko-ins-" + instanceID
	instanceCheckID = "reko-ins-chk-" + instanceID
	return
}

func registerInstance(bind string) (err error) {
	var addr *net.TCPAddr
	if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
		return
	}

	if err = cclient.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Name: instanceServiceName,
		ID:   instanceServiceID,
		Port: addr.Port,
		Check: &consul.AgentServiceCheck{
			CheckID: instanceCheckID,
			Name:    "Internal Alive Check",
			TTL:     "10s",
		},
	}); err != nil {
		return
	}
	return
}

func deregisterInstance() error {
	return cclient.Agent().ServiceDeregister(instanceServiceID)
}

func notifyInstanceRunning() error {
	return cclient.Agent().PassTTL(instanceCheckID, "RUNNING")
}

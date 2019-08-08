package main

import (
	"flag"
	consul "github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	optDeregister bool

	envBind = os.Getenv("BIND")

	cclient *consul.Client
)

func init() {
	if len(envBind) == 0 {
		envBind = "0.0.0.0:9001"
	}
}

func exit(err *error) {
	if *err != nil {
		log.Printf("exited with error: %s", (*err).Error())
		os.Exit(1)
	}
}

func main() {
	var err error
	defer exit(&err)

	flag.BoolVar(&optDeregister, "deregister", false, "one shot run to deregister self from consul")
	flag.Parse()

	if err = ensureInstanceID(); err != nil {
		return
	}

	if cclient, err = consul.NewClient(consul.DefaultConfig()); err != nil {
		return
	}

	if optDeregister {
		err = deregisterInstance()
		return
	}

	if err = registerInstance(envBind); err != nil {
		return
	}

	go watchdog()

	_ = http.ListenAndServe(envBind, NewHandler(cclient))
}

func watchdog() {
	tk := time.NewTicker(time.Second * 5)
	defer tk.Stop()

	for {
		<-tk.C
		if err := notifyInstanceRunning(); err != nil {
			log.Printf("failed to notify instance running: %s", err.Error())
		}
	}
}

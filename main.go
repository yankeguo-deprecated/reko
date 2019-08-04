package main

import (
	consul "github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"os"
)

var (
	optBind = os.Getenv("BIND")
)

func init() {
	if len(optBind) == 0 {
		optBind = "127.0.0.1:9001"
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

	var client *consul.Client
	if client, err = consul.NewClient(consul.DefaultConfig()); err != nil {
		return
	}

	_ = http.ListenAndServe(optBind, &Handler{Client: client})
}

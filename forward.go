package main

import (
	"fmt"
	"net/http"
)

func Forward(hosts []string, rw http.ResponseWriter, r *http.Request) {
	// TODO: not implemented
	rw.WriteHeader(http.StatusBadGateway)
	_, _ = rw.Write([]byte(fmt.Sprintf("reko: not implemented, %v, %s", hosts, r.URL.String())))
}

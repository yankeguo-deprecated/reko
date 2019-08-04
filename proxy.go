package main

import "net/http"

func Proxy(key string, hosts []string, rw http.ResponseWriter, r *http.Request) {
	// TODO: not implemented
	rw.WriteHeader(http.StatusBadGateway)
	_, _ = rw.Write([]byte("reko: not implements"))
}

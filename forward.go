package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

type Proxy struct {
	Hosts        []string
	Cursor       int
	ReverseProxy httputil.ReverseProxy
}

func (p *Proxy) Director(r *http.Request) {
	r.URL.Scheme = "http"
	r.URL.Host = p.Hosts[p.Cursor]
}

func (p *Proxy) ErrorHandler(rw http.ResponseWriter, r *http.Request, err error) {
	log.Printf("upstream %s failed", p.Hosts[p.Cursor])
	p.Cursor++
	if p.Cursor < len(p.Hosts) {
		p.ReverseProxy.ServeHTTP(rw, r)
		return
	}
	rw.WriteHeader(http.StatusServiceUnavailable)
	_, _ = rw.Write([]byte(fmt.Sprintf("reko: all services failed")))
}

func NewProxy(hosts []string) httputil.ReverseProxy {
	p := &Proxy{Hosts: hosts}
	p.ReverseProxy.Director = p.Director
	p.ReverseProxy.ErrorHandler = p.ErrorHandler
	return p.ReverseProxy
}

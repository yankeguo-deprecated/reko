package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

type Upstream struct {
	Name string
	Host string
}

type Proxy struct {
	Upstreams    []Upstream
	Cursor       int
	ReverseProxy httputil.ReverseProxy
}

func (p *Proxy) Director(r *http.Request) {
	r.URL.Scheme = "http"
	r.URL.Host = p.Upstreams[p.Cursor].Host
}

func (p *Proxy) ErrorHandler(rw http.ResponseWriter, r *http.Request, err error) {
	log.Printf("upstream %+v failed: %s", p.Upstreams[p.Cursor], err.Error())
	p.Cursor++
	if p.Cursor < len(p.Upstreams) {
		p.ReverseProxy.ServeHTTP(rw, r)
		return
	}
	rw.WriteHeader(http.StatusServiceUnavailable)
	_, _ = rw.Write([]byte(fmt.Sprintf("reko: all upstreams failed")))
}

func (p *Proxy) ModifyResponse(res *http.Response) error {
	res.Header.Set("X-Reko-Upstream", p.Upstreams[p.Cursor].Name)
	return nil
}

func NewProxy(upstreams []Upstream) httputil.ReverseProxy {
	p := &Proxy{Upstreams: upstreams}
	p.ReverseProxy.Director = p.Director
	p.ReverseProxy.ErrorHandler = p.ErrorHandler
	p.ReverseProxy.ModifyResponse = p.ModifyResponse
	return p.ReverseProxy
}

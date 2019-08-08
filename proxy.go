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

type reverseProxy struct {
	*httputil.ReverseProxy
	Upstreams []Upstream
	Cursor    int
}

func (p *reverseProxy) director(r *http.Request) {
	r.URL.Scheme = "http"
	r.URL.Host = p.Upstreams[p.Cursor].Host
}

func (p *reverseProxy) errorHandler(rw http.ResponseWriter, r *http.Request, err error) {
	log.Printf("upstream %+v failed: %s", p.Upstreams[p.Cursor], err.Error())
	p.Cursor++
	if p.Cursor < len(p.Upstreams) {
		p.ReverseProxy.ServeHTTP(rw, r)
		return
	}
	rw.WriteHeader(http.StatusServiceUnavailable)
	_, _ = rw.Write([]byte(fmt.Sprintf("reko: all upstreams failed")))
}

func (p *reverseProxy) modifyResponse(res *http.Response) error {
	res.Header.Set("X-Reko-Upstream", p.Upstreams[p.Cursor].Name)
	return nil
}

func NewProxy(upstreams []Upstream) http.Handler {
	p := &reverseProxy{
		Upstreams: upstreams,
	}
	p.ReverseProxy = &httputil.ReverseProxy{
		Director:       p.director,
		ModifyResponse: p.modifyResponse,
		ErrorHandler:   p.errorHandler,
	}
	return p
}

package main

import (
	"errors"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"net/http"
	"sync/atomic"
)

type Handler struct {
	Client *consul.Client

	RR map[string]*uint64
}

func NewHandler(client *consul.Client) *Handler {
	return &Handler{
		Client: client,
		RR:     map[string]*uint64{},
	}
}

func (h *Handler) NextRR(key string) uint64 {
	// no locking, allow inconsistency for better perf
	rr := h.RR[key]
	if rr == nil {
		rr = new(uint64)
		h.RR[key] = rr
	}
	// increase
	return atomic.AddUint64(rr, 1)
}

func (h *Handler) Rotate(key string, hosts []string) []string {
	rr := h.NextRR(key)
	cr := rr % uint64(len(hosts))
	return append(hosts[cr:], hosts[0:cr]...)
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var err error
	defer func(err *error) {
		if *err != nil {
			rw.WriteHeader(http.StatusBadGateway)
			_, _ = rw.Write([]byte(fmt.Sprintf("reko: %s", (*err).Error())))
		}
	}(&err)

	var q ServiceQuery
	if q, err = ExtractServiceQuery(r.URL); err != nil {
		return
	}

	var hosts []string
	if hosts, err = q.Resolve(h.Client); err != nil {
		return
	}

	if len(hosts) == 0 {
		err = errors.New("no services available")
		return
	}

	// rotate hosts
	hosts = h.Rotate(q.Raw, hosts)

	// proxy
	Forward(hosts, rw, r)
}

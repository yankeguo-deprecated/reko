package main

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Handler struct {
	Client *consul.Client
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var err error
	defer func(err *error) {
		if *err != nil {
			rw.WriteHeader(http.StatusBadGateway)
			_, _ = rw.Write([]byte(fmt.Sprintf("reko: %s", (*err).Error())))
		}
	}(&err)
	var sName string
	var sTags []string
	if sName, sTags, err = normalizeIncomingURL(r.URL); err != nil {
		return
	}
	var ss []*consul.CatalogService
	if ss, _, err = h.Client.Catalog().ServiceMultipleTags(sName, sTags, &consul.QueryOptions{AllowStale: true}); err != nil {
		return
	}
	if len(ss) == 0 {
		if len(sTags) > 0 {
			err = fmt.Errorf("service '%s' with tags %v not found", sName, sTags)
		} else {
			err = fmt.Errorf("service '%s' not found", sName)
		}
		return
	}
	hosts := make([]string, 0, len(ss))
	for _, s := range ss {
		if len(s.ServiceAddress) > 0 {
			hosts = append(hosts, s.ServiceAddress+":"+strconv.Itoa(s.ServicePort))
		} else {
			hosts = append(hosts, s.Address+":"+strconv.Itoa(s.ServicePort))
		}
	}
	Proxy(sName+"-"+strings.Join(sTags, "-"), hosts, rw, r)
}

func normalizeIncomingURL(u *url.URL) (sName string, sTags []string, err error) {
	cs := strings.Split(u.Path, "/")
	var queryIdx int
	var query string
	for i, c := range cs {
		if len(c) > 0 {
			queryIdx, query = i, c
			break
		}
	}
	if len(query) == 0 {
		err = fmt.Errorf("missing service query: %s", u.Path)
		return
	}
	queryParts := strings.Split(query, ",")
	if len(queryParts) == 0 {
		err = fmt.Errorf("invalid service query: %s", query)
		return
	}
	for _, c := range queryParts {
		if len(c) == 0 {
			err = fmt.Errorf("invalid service query: %s", query)
			return
		}
	}
	sName = queryParts[0]
	sTags = queryParts[1:]
	cs[queryIdx] = sName
	u.Path = strings.Join(cs, "/")
	return
}

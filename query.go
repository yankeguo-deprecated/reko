package main

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"net/url"
	"strconv"
	"strings"
)

type ServiceQuery struct {
	Raw  string
	Name string
	ID   string
	Tags []string
}

func ExtractServiceQuery(u *url.URL) (q ServiceQuery, err error) {
	parts := strings.Split(u.Path, "/")
	var idx int
	var raw string
	for i, p := range parts {
		if len(p) > 0 {
			idx, raw = i, p
			break
		}
	}
	if len(raw) == 0 {
		err = fmt.Errorf("missing service query: %s", u.Path)
		return
	}

	q.Raw = raw

	if tags := strings.Split(raw, ":"); len(tags) == 2 {
		q.Name = tags[0]
		q.Tags = strings.Split(tags[1], ",")
	} else if ids := strings.Split(raw, "@"); len(ids) == 2 {
		q.Name = ids[0]
		q.ID = ids[1]
	} else {
		q.Name = raw
	}

	u.Path = strings.Join(append(parts[0:idx], parts[idx+1:]...), "/")
	if len(u.Path) == 0 {
		u.Path = "/"
	}
	return
}

func (q ServiceQuery) Resolve(client *consul.Client) (hosts []string, err error) {
	var ret []*consul.CatalogService
	if ret, _, err = client.Catalog().ServiceMultipleTags(q.Name, q.Tags, &consul.QueryOptions{AllowStale: true}); err != nil {
		return
	}
	for _, s := range ret {
		if len(q.ID) == 0 || s.ServiceID == q.ID {
			if len(s.ServiceAddress) > 0 {
				hosts = append(hosts, s.ServiceAddress+":"+strconv.Itoa(s.ServicePort))
			} else {
				hosts = append(hosts, s.Address+":"+strconv.Itoa(s.ServicePort))
			}
		}
	}
	return
}

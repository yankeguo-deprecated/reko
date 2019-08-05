package main

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestExtractServiceQuery(t *testing.T) {
	var u *url.URL
	var q ServiceQuery
	var err error

	u, _ = url.Parse("http://127.0.0.1:9001/demo-service:canary,test/hello/world?key=val")
	q, err = ExtractServiceQuery(u)
	assert.NoError(t, err, "no error")
	assert.Equal(t, "demo-service", q.Name, "service name")
	assert.Equal(t, []string{"canary", "test"}, q.Tags, "service tags")
	assert.Equal(t, "/hello/world", u.Path, "twisted path")
	assert.Equal(t, "http://127.0.0.1:9001/hello/world?key=val", u.String(), "url string")

	u, _ = url.Parse("http://127.0.0.1:9001/demo-service:canary,test//hello/world/?key=val")
	q, err = ExtractServiceQuery(u)
	assert.NoError(t, err, "no error")
	assert.Equal(t, "demo-service", q.Name, "service name")
	assert.Equal(t, []string{"canary", "test"}, q.Tags, "service tags")
	assert.Equal(t, "//hello/world/", u.Path, "twisted path")
	assert.Equal(t, "http://127.0.0.1:9001//hello/world/?key=val", u.String(), "url string")

	u, _ = url.Parse("http://127.0.0.1:9001/demo-service@service-1//hello/world/?key=val")
	q, err = ExtractServiceQuery(u)
	assert.NoError(t, err, "no error")
	assert.Equal(t, "demo-service", q.Name, "service name")
	assert.Equal(t, "service-1", q.ID, "service name")
	assert.Equal(t, "//hello/world/", u.Path, "twisted path")
	assert.Equal(t, "http://127.0.0.1:9001//hello/world/?key=val", u.String(), "url string")
}

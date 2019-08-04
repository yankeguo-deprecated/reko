package main

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestNormalizeIncomingURL(t *testing.T) {
	var u *url.URL
	var sName string
	var sTags []string
	var err error

	u, _ = url.Parse("http://127.0.0.1:9001/demo-service,canary,test/hello/world?key=val")
	sName, sTags, err = normalizeIncomingURL(u)
	assert.NoError(t, err, "no error")
	assert.Equal(t, "demo-service", sName, "service name")
	assert.Equal(t, []string{"canary", "test"}, sTags, "service tags")
	assert.Equal(t, "/demo-service/hello/world", u.Path, "twisted path")
	assert.Equal(t, "http://127.0.0.1:9001/demo-service/hello/world?key=val", u.String(), "url string")

	u, _ = url.Parse("http://127.0.0.1:9001/demo-service,canary,test//hello/world/?key=val")
	sName, sTags, err = normalizeIncomingURL(u)
	assert.NoError(t, err, "no error")
	assert.Equal(t, "demo-service", sName, "service name")
	assert.Equal(t, []string{"canary", "test"}, sTags, "service tags")
	assert.Equal(t, "/demo-service//hello/world/", u.Path, "twisted path")
	assert.Equal(t, "http://127.0.0.1:9001/demo-service//hello/world/?key=val", u.String(), "url string")
}

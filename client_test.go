package ztcentral

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
)

func TestHeaders(t *testing.T) {
	type testinfo struct {
		pass       bool
		clientFunc func() *Client
		matchFunc  func(http.Header) bool
	}

	table := map[string]testinfo{
		"plain user-agent positive": {
			pass: true,
			clientFunc: func() *Client {
				return NewClient("test")
			},
			matchFunc: func(h http.Header) bool {
				return h.Get("User-Agent") == userAgent
			},
		},
		"plain user-agent negative": {
			pass: false,
			clientFunc: func() *Client {
				c := NewClient("test")
				c.SetUserAgent("stuff")
				return c
			},
			matchFunc: func(h http.Header) bool {
				return h.Get("User-Agent") == userAgent
			},
		},
		"modified user-agent positive": {
			pass: true,
			clientFunc: func() *Client {
				c := NewClient("test")
				c.SetUserAgent("stuff")
				return c
			},
			matchFunc: func(h http.Header) bool {
				return h.Get("User-Agent") == fmt.Sprintf("%s (stuff)", userAgent)
			},
		},
	}

	for name, info := range table {
		c := info.clientFunc()
		req, err := retryablehttp.NewRequest("GET", "http://localhost", nil)
		if err != nil {
			t.Fatalf("wtf, %s?", name)
		}

		res := info.matchFunc(c.prepareRequest(context.Background(), req).Header)
		if res != info.pass {
			t.Fatalf("%q: result was expected to be %v but was %v", name, info.pass, res)
		}
	}
}

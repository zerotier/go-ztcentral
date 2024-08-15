// Copyright (c) 2021, ZeroTier, Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Package ztcentral provides an API for interacting with ZeroTier Central (https://my.zerotier.com)
package ztcentral

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/zerotier/go-ztcentral/pkg/spec"
)

var (
	// BaseURLV1 is the URL to ZeroTier Central
	BaseURLV1 = "https://my.zerotier.com/api"
)

var userAgent = fmt.Sprintf("go-ztcentral/%s", Version)

// Client is the zerotier central client.
type Client struct {
	specClient *spec.Client
	httpClient *http.Client

	apiKey    string
	userAgent string
	limits    RateLimitHeaders
}

type RateLimitHeaders struct {
	Limit     int       `json:"x-ratelimit-limit"`
	Remaining int       `json:"x-ratelimit-remaining"`
	ResetTime ResetTime `json:"x-ratelimit-reset"`
}

type ResetTime string

func (r ResetTime) String() string {
	return string(r)
}

func (r ResetTime) Time() time.Time {
	d, _ := time.ParseDuration(string(r))
	return time.Now().Add(d)
}

func newRateLimitHeaders(h http.Header) RateLimitHeaders {
	// Central Sends Headers:
	// X-Ratelimit-Limit: 20
	// X-Ratelimit-Remaining: 14
	// X-Ratelimit-Reset: Tue, 06 Aug 2024 20:53:48 UTC

	limitReq, _ := strconv.Atoi(h.Get("x-ratelimit-limit"))
	remainingReq, _ := strconv.Atoi(h.Get("x-ratelimit-remaining"))
	resetTime := ResetTime(h.Get("x-ratelimit-reset"))

	return RateLimitHeaders{
		Limit:     limitReq,
		Remaining: remainingReq,
		ResetTime: resetTime,
	}
}

// ErrStatus is returned when the response code is not 200.
var ErrStatus = errors.New("status code was not 200")

// NewClient creates a client.
// key is an API key for your ZeroTier Central that you can generate after login.
// It returns a fully initialized client.
func NewClient(key string) (*Client, error) {
	c := &Client{
		apiKey:    key,
		userAgent: userAgent,
	}

	c.httpClient = &http.Client{Transport: c}

	var err error
	c.specClient, err = spec.NewClient(BaseURLV1, spec.WithHTTPClient(c.httpClient))
	if err != nil {
		return nil, err
	}

	return c, nil
}

// SetUserAgent appends a custom user agent to the existing one, allowing
// customization of it. It will be present where browsers typically put
// extension data, such as Mozilla/5.0 (Firefox 80). The "Firefox 80" section
// will be replaced with what is provided here.
//
// While this code is intended to be used by third party code, be advised that
// this does nothing but identify your product to the zerotier network. Nothing
// you put here should change the client behavior.
func (c *Client) SetUserAgent(ua string) {
	c.userAgent = fmt.Sprintf("%s (%s)", c.userAgent, ua)
}

// RoundTrip conforms the client to http.RoundTrip
func (c *Client) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.apiKey))

	if c.limits.Limit != 0 && c.limits.Remaining < c.limits.Limit {
		diff := time.Until(time.Now().Add(time.Duration(c.limits.Limit-c.limits.Remaining) * 10 * time.Millisecond))
		time.Sleep(diff)
	}

	resp, err := http.DefaultTransport.RoundTrip(req)

	a := newRateLimitHeaders(resp.Header)

	c.limits = a

	return resp, err
}

func (c *Client) decode(resp *http.Response, i interface{}) error {
	if resp.StatusCode != 200 {
		return fmt.Errorf("Status code %v: %w", resp.StatusCode, ErrStatus)
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(i)
}

func (c *Client) decomposeStruct(i interface{}) (map[string]interface{}, error) {
	res, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	return m, json.Unmarshal(res, &m)
}

// User gets the user of the client API token, via the /status endpoint.
//
// For the full status, see Status()
func (c *Client) User(ctx context.Context) (*spec.User, error) {
	res, err := c.Status(ctx)
	if err != nil {
		return nil, err
	}

	return res.User, nil
}

// Status returns the full status response of the /status API endpoint, which
// contains various bits of information about the client's account.
//
// For just the User information, see User().
func (c *Client) Status(ctx context.Context) (*spec.Status, error) {
	res := &spec.Status{}
	resp, err := c.specClient.GetStatus(ctx)
	if err != nil {
		return nil, err
	}

	return res, c.decode(resp, res)
}

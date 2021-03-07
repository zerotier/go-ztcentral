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

	"github.com/hashicorp/go-retryablehttp"
)

const (
	// BaseURLV1 is the URL to ZeroTier Central
	BaseURLV1 = "https://my.zerotier.com/api"
)

var userAgent = fmt.Sprintf("go-ztcentral/%s", Version)

// Client is the zerotier central client.
type Client struct {
	BaseURL    string
	HTTPClient *retryablehttp.Client

	apiKey    string
	userAgent string
}

// NewClient creates a client.
// key is an API key for your ZeroTier Central that you can generate after login.
// It returns a fully initialized client.
func NewClient(key string) *Client {
	c := &Client{
		BaseURL:    BaseURLV1,
		HTTPClient: retryablehttp.NewClient(),

		apiKey:    key,
		userAgent: userAgent,
	}

	return c
}

// ErrorResponse is the response to an error
type ErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int
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

func (c *Client) prepareRequest(ctx context.Context, req *retryablehttp.Request) *retryablehttp.Request {
	req = req.WithContext(ctx)

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.apiKey))

	return req
}

func (c *Client) sendRequest(ctx context.Context, req *retryablehttp.Request, v interface{}) error {
	res, err := c.HTTPClient.Do(c.prepareRequest(ctx, req))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if v != nil {
		if err = json.NewDecoder(res.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}

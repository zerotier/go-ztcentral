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

package ztcentral

import (
	"context"
	"errors"
	"fmt"

	"github.com/zerotier/go-ztcentral/pkg/spec"
)

// CreateAPIToken creates an API token with the secret you desire in Central.
func (c *Client) CreateAPIToken(ctx context.Context, userID, name, token string) error {
	if len(token) < 32 {
		return errors.New("token must be a minimum of 32 characters")
	}

	resp, err := c.specClient.AddAPIToken(ctx, userID, spec.AddAPITokenJSONRequestBody{
		Token:     &token,
		TokenName: &name,
	})

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Response was not 200; was %d", resp.StatusCode)
	}

	return nil
}

// DeleteAPIToken removes an API token from the list of available tokens.
func (c *Client) DeleteAPIToken(ctx context.Context, userID, name string) error {
	resp, err := c.specClient.DeleteAPIToken(ctx, userID, name)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Response was not 200; was %d", resp.StatusCode)
	}

	return nil
}

// RandomToken fetches an API-compatible token that can be fed to CreateAPIToken.
func (c *Client) RandomToken(ctx context.Context) (string, error) {
	res := spec.RandomToken{}

	resp, err := c.specClient.GetRandomToken(ctx)
	if err != nil {
		return "", err
	}

	if err := c.decode(resp, &res); err != nil {
		return "", err
	}

	return *res.Token, nil
}

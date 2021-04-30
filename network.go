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
	"fmt"

	"github.com/zerotier/go-ztcentral/pkg/spec"
)

// GetNetworks returns the list of your available networks
func (c *Client) GetNetworks(ctx context.Context) ([]spec.Network, error) {
	var res []spec.Network
	resp, err := c.specClient.GetNetworkList(ctx)
	if err != nil {
		return res, err
	}

	return res, c.decode(resp, &res)
}

// GetNetwork returns an individual network specified by networkID
func (c *Client) GetNetwork(ctx context.Context, networkID string) (spec.Network, error) {
	var res spec.Network

	resp, err := c.specClient.GetNetworkByID(ctx, networkID)
	if err != nil {
		return res, err
	}

	return res, c.decode(resp, &res)
}

func (c *Client) UpdateNetwork(ctx context.Context, id string, network spec.Network) (spec.Network, error) {
	var res spec.Network

	resp, err := c.specClient.UpdateNetwork(ctx, id, spec.UpdateNetworkJSONRequestBody(network))
	if err != nil {
		return res, err
	}

	return res, c.decode(resp, &res)
}

func (c *Client) UpdateNetworkRules(ctx context.Context, id, source string) (string, error) {
	net, err := c.UpdateNetwork(ctx, id, spec.Network{Id: &id, RulesSource: &source})
	if err != nil {
		return "", err
	}

	if net.RulesSource == nil {
		return "", nil
	}

	return *net.RulesSource, nil
}

func (c *Client) NewNetwork(ctx context.Context, name string, n spec.Network) (spec.Network, error) {
	if n.Config != nil {
		n.Config.Name = &name
	} else {
		n.Config = &spec.NetworkConfig{Name: &name}
	}

	var newnet spec.Network

	net, err := c.decomposeStruct(n)
	if err != nil {
		return newnet, err
	}

	resp, err := c.specClient.NewNetwork(ctx, net)
	if err != nil {
		return newnet, err
	}

	return newnet, c.decode(resp, &newnet)
}

func (c *Client) DeleteNetwork(ctx context.Context, networkID string) error {
	resp, err := c.specClient.DeleteNetwork(ctx, networkID)
	if err != nil {
		return err
	}

	resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Status was %v: %w", resp.StatusCode, ErrStatus)
	}

	return nil
}

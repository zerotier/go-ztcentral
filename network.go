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
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-retryablehttp"
)

type Network struct {
	ID                    string                 `json:"id"`
	Type                  string                 `json:"type"`
	Clock                 int64                  `json:"clock"`
	Config                NetworkConfig          `json:"config"`
	Description           string                 `json:"description"`
	RulesSource           string                 `json:"rulesSource"`
	Permissions           NetworkPermissionsMap  `json:"permissions"`
	OwnerID               string                 `json:"ownerId"`
	OnlineMemberCount     int                    `json:"onlineMemberCount"`
	AuthorizedMemberCount int                    `json:"authorizedMemberCount"`
	TotalMemberCount      int                    `json:"totalMemberCount"`
	CapabilitiesByName    map[string]interface{} `json:"capabilitiesByName"`
	TagsByName            map[string]interface{} `json:"tagsByName"`
	UI                    map[string]interface{} `json:"ui"`
}

type NetworkConfig struct {
	ActiveMemberCount int                    `json:"activeMemberCount,omitempty"`
	CreationTime      int64                  `json:"creationTime"`
	Capabilities      []interface{}          `json:"capabilities"`
	EnableBroadcast   bool                   `json:"enableBroadcast"`
	ID                string                 `json:"id,omitempty"`
	IPAssignmentPool  []IPRange              `json:"ipAssignmentPools,omitempty"`
	LastModified      int64                  `json:"lastModified"`
	MTU               int                    `json:"mtu,omitempty"`
	MulticastLimit    int                    `json:"multicastLimit"`
	Name              string                 `json:"name"`
	Private           bool                   `json:"private"`
	RemoteTraceLevel  int                    `json:"remoteTraceLevel"`
	RemoteTraceTarget *string                `json:"remoteTraceTarget"`
	Revision          uint64                 `json:"revision,omitempty"`
	Routes            []Route                `json:"routes"`
	Rules             []interface{}          `json:"rules"`
	Tags              []interface{}          `json:"tags"`
	IPV4AssignMode    map[string]interface{} `json:"v4AssignMode"`
	IPV6AssignMode    map[string]interface{} `json:"v6AssignMode"`
	DNS               *NetworkDNS            `json:"dns,omitempty"`
}

type NetworkPermissions struct {
	Authorize bool `json:"a"`
	Delete    bool `json:"d"`
	Modify    bool `json:"m"`
	Read      bool `json:"r"`
}

type NetworkPermissionsMap map[string]NetworkPermissions

type IPRange struct {
	Start string `json:"ipRangeStart"`
	End   string `json:"ipRangeEnd"`
}

type Route struct {
	Target string `json:"target"`
	Via    string `json:"via,omitempty"`
}

type NetworkDNS struct {
	Domain  string   `json:"domain"`
	Servers []string `json:"servers"`
}

type NetworkList []Network

func (c *Client) GetNetworks(ctx context.Context) (NetworkList, error) {
	req, err := retryablehttp.NewRequest("GET", fmt.Sprintf("%s/network", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	res := make(NetworkList, 0)
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetNetwork(ctx context.Context, networkID string) (*Network, error) {
	req, err := retryablehttp.NewRequest("GET", fmt.Sprintf("%s/network/%s", c.BaseURL, networkID), nil)
	if err != nil {
		return nil, err
	}

	res := Network{}
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) UpdateNetwork(ctx context.Context, network *Network) (*Network, error) {
	reqBody, err := json.Marshal(network)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/network/%s", c.BaseURL, network.ID), reqBody)
	if err != nil {
		return nil, err
	}

	res := Network{}
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) NewNetwork(ctx context.Context, name string) (*Network, error) {
	n := Network{
		Config: NetworkConfig{
			Name: name,
		},
	}

	body, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/network", c.BaseURL), body)
	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(ctx, req, &n); err != nil {
		return nil, err
	}

	return &n, nil
}

func (c *Client) DeleteNetwork(ctx context.Context, networkID string) error {
	req, err := retryablehttp.NewRequest("DELETE", fmt.Sprintf("%s/network/%s", c.BaseURL, networkID), nil)
	if err != nil {
		return err
	}

	return c.sendRequest(ctx, req, nil)
}

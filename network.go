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

// Network represents a ZeroTier Network at https://my.zerotier.com
type Network struct {

	// The 16-digit hex Network ID [read only]
	ID string `json:"id"`

	// The type of the JSON object [read only]
	Type string `json:"type"`

	// Current server system clock [read only]
	Clock int64 `json:"clock"`

	// NetworkConfig sub-object
	Config NetworkConfig `json:"config"`

	// Description of the network
	Description string `json:"description"`

	// Network ruels engine source code
	RulesSource string `json:"rulesSource"`

	// Network editing access permissions
	Permissions NetworkPermissionsMap `json:"permissions"`

	// ZeroTeir Central user ID of the owner of the network [read only]
	OwnerID string `json:"ownerId"`

	// Count of current online members [read only]
	OnlineMemberCount int `json:"onlineMemberCount"`

	// Count of members authorized on the network [read only]
	AuthorizedMemberCount int `json:"authorizedMemberCount"`

	// Total number of members with access or requesting access to the network [read only]
	TotalMemberCount int `json:"totalMemberCount"`

	// Capabilties defined in rule set by name
	CapabilitiesByName map[string]interface{} `json:"capabilitiesByName"`

	// Tags defined in rule set by name
	TagsByName map[string]interface{} `json:"tagsByName"`
}

// NetworkConfig object represents individual configuration options on the network
type NetworkConfig struct {
	// Time of network creation
	CreationTime int64 `json:"creationTime"`

	// Array of capabilities available on this network (see https://www.zerotier.com/manual/#3)
	Capabilities []interface{} `json:"capabilities"`

	// Whether or not Broadcast packets are allowed on the network
	EnableBroadcast bool `json:"enableBroadcast"`

	// 16 digit hexidecimal Network ID
	ID string `json:"id,omitempty"`

	// Array of available IPRangees to use for address assignment
	IPAssignmentPool []IPRange `json:"ipAssignmentPools,omitempty"`

	// Time of last network modification
	LastModified int64 `json:"lastModified"`

	// MTU of the virtual network (default 2800)
	MTU int `json:"mtu,omitempty"`

	// Maximum number of recipients per multicast or broadcast
	//
	// Warning: Setting this to 0 will disable IPv4 communication on your network!
	MulticastLimit int `json:"multicastLimit"`

	// Name of the network
	Name string `json:"name"`

	// Whether this is a private or public networks
	//
	// If false, members do not require authorization to join the network.
	// They will be auto-accepted & authorized.
	Private bool `json:"private"`

	// Network configuration revision number [read only]
	Revision uint64 `json:"revision,omitempty"`

	// Array of Routes published to members
	Routes []Route `json:"routes"`

	// Network base rules (see https://www.zerotier.com/manual/#3)
	Rules []interface{} `json:"rules"`

	// Array of tags available on the network (see https://www.zerotier.com/manual/#3)
	Tags []interface{} `json:"tags"`

	// IPv4 address assignment modes
	IPV4AssignMode `json:"v4AssignMode"`

	// IPv6 address assignment modes
	IPV6AssignMode `json:"v6AssignMode"`

	// Network DNS configuration
	DNS *NetworkDNS `json:"dns,omitempty"`
}

// NetworkPermissions holds the 4 different permissions settable on a network
type NetworkPermissions struct {
	Authorize bool `json:"a"`
	Delete    bool `json:"d"`
	Modify    bool `json:"m"`
	Read      bool `json:"r"`
}

// NetworkPermissionsMap is a map of a ZeroTier Central User ID to NetworkPermissions
type NetworkPermissionsMap map[string]NetworkPermissions

// IPRange represents a range of IP addresses from start to end.
// Can be either IPv4 or IPv6
type IPRange struct {
	Start string `json:"ipRangeStart"`
	End   string `json:"ipRangeEnd"`
}

// Route is a network route uset to specify a managed route published to clients on a network
type Route struct {
	// CIDR of the network target for the route
	Target string `json:"target"`

	// Optional IP address to route the Target via
	Via string `json:"via,omitempty"`
}

// NetworkDNS holds DNS information published to a network
type NetworkDNS struct {
	// Search domain
	Domain string `json:"domain"`

	// Array of up to 4 IP addresses that will be DNS servers
	Servers []string `json:"servers"`
}

// NetworkList is an array of Network structs
type NetworkList []Network

// IPV4AssignMode allows enabling disabling of IPv4 address assignment modes.
// Currently there is only one.
type IPV4AssignMode struct {
	// Network controller assigns IP addresses from the IPv4 auto-assign range.
	// If false and/or no Auto Assign Range is specified, the user must manually
	// specify all IP addresses for network members
	ZeroTier bool `json:"zt"`
}

// IPV6AssignMode allows enabling and disabling of IPv6 address assignment modes.
type IPV6AssignMode struct {
	// 6PLANE assigns every host on a ZeroTier network an IPv6 address within a
	// private /40 network. The 8-bit fc prefix indicates a private IPv6 network with
	// an "experimental" assignment scheme (not important here), while the remaining
	// 32 bits are computed by XORing the upper and lower 32 bits of the network ID
	// together. This yields a unique deterministic prefix for every ZeroTier virtual
	// network.
	//
	// Inside this network the controller will hand out IPv6 unicast addresses like
	// fcf9:b03a:1289:e92c:eee5::1 to every participant. Look closely and you'll see
	// our /40 followed by another 40 bits: 89:e92c:eee5. This is the 40-bit ZeroTier
	// address of the host.
	ZT6Plane bool `json:"6plane"`

	// Assign a unique /128 for each device
	RFC4193 bool `json:"rfc4193"`

	// Network controller assigns IP addresses from the IPv6 auto-assign range.
	// If false and/or no Auto Assign Range is specified, the user must manually
	// specify all IP addresses for network members
	ZeroTier bool `json:"zt"`
}

// GetNetworks returns the list of your available networks
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

// GetNetwork returns an individual network specified by networkID
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

func (c *Client) NewNetwork(ctx context.Context, name string, n *Network) (*Network, error) {
	if n != nil {
		n.Config.Name = name
	} else {
		n = &Network{Config: NetworkConfig{Name: name}}
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

	return n, nil
}

func (c *Client) DeleteNetwork(ctx context.Context, networkID string) error {
	req, err := retryablehttp.NewRequest("DELETE", fmt.Sprintf("%s/network/%s", c.BaseURL, networkID), nil)
	if err != nil {
		return err
	}

	return c.sendRequest(ctx, req, nil)
}

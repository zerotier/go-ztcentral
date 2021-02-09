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

// Member represents a member of a ZeroTier network
type Member struct {
	// (Deprecated) ID is a concatenation of NetworkID-MemberID [read only]
	ID string `json:"id,omitempty"`

	// The type of the JSON object [read only]
	Type string `json:"type,omitempty"`

	// Current server system clock [read only]
	Clock int64 `json:"clock,omitempty"`

	// ID of the Network the member belongs to.
	NetworkID string `json:"networkId,omitempty"`

	// ID of the member node.  This is the 10 digit identifier that identifies a ZeroTier node. [read only]
	MemberID string `json:"nodeId,omitempty"`

	// (Deprecated) ID of the network controller [read only]
	ControllerID string `json:"controllerId,omitempty"`

	// Whether or not the member is hidden in the uI
	Hidden bool `json:"hidden"`

	// Name of the network member
	Name string `json:"name,omitempty"`

	// Whether or not the member is currently online [read only]
	Online bool `json:"online"`

	// Description of the member
	Description string `json:"description,omitempty"`

	// MemberConfig sub struct
	Config MemberConfig `json:"config"`

	// Last time the member was connected to the controller. [read only]
	LastOnline int64 `json:"lastOnline,omitempty"`

	// IP address the member last contacted the controller from [read only]
	PhysicalAddress *string `json:"physicalAddress,omitempty"`

	// Version string of the client (ex: 1.6.3) [read only]
	ClientVersion string `json:"clientVersion,omitempty"`

	// Protocol version of the client [read only]
	ProtocolVersion int `json:"protocolVersion"`

	// Whether the client is new enough to support the rules engine [read only]
	SupportsRulesEngine bool `json:"supportsRulesEngine,omitempty"`
}

// MemberConfig represents individual configuration options of a network member
type MemberConfig struct {
	// Allow the member to be a Bridge on a network
	ActiveBridge bool `json:"activeBridge"`

	// Whether or not the member is authorized to communicate on a network
	Authorized bool `json:"authorized"`

	// Array of IDs of capabilities assigned to this member
	Capabilities []uint `json:"capabilities,omitempty"`

	// Time the member was created or first tried to join the network [read only]
	CreationTime int64 `json:"creationTime"`

	// ID of the member node.  This is the 10 digit identifier that identifies a ZeroTier node. [read only]
	MemberID string `json:"id,omitempty"`

	// Public Key of the member's Identity [read only]
	Identity string `json:"identity,omitempty"`

	// List of IP addresses assigned to the member
	IPAssignments []string `json:"ipAssignments,omitempty"`

	// Time the member was authorized on the network
	LastAuthorizedTime int64 `json:"lastAuthorizedTime"`

	// Time the member was last deauthorized on the network.
	LastDeauthorizedTime int64 `json:"lastDeauthorizedTime"`

	// Exempt this member from the IP auto assignment pool on a Network
	NoAutoAssignIPs bool `json:"noAutoAssignIps"`

	// Member record revision count [read only]
	Revision uint64 `json:"revision"`

	// Array of tuples of tag ID, tag value
	Tags [][2]uint `json:"tags,omitempty"`

	// Major version of the client
	VersionMajor int `json:"vMajor"`

	// Minor version of the clinet
	VersionMinor int `json:"vMinor"`

	// Revision number of client
	VersionRev int `json:"vRev"`

	// Protocol version number
	VersionProtocol int `json:"vProto"`
}

// MemberList is an array of Member structs
type MemberList []Member

func (c *Client) GetMembers(ctx context.Context, networkID string) (MemberList, error) {
	req, err := retryablehttp.NewRequest("GET", fmt.Sprintf("%s/network/%s/member", c.BaseURL, networkID), nil)
	if err != nil {
		return nil, err
	}

	res := make(MemberList, 0)
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetMember(ctx context.Context, networkID, memberID string) (*Member, error) {
	req, err := retryablehttp.NewRequest("GET", fmt.Sprintf("%s/network/%s/member/%s", c.BaseURL, networkID, memberID), nil)
	if err != nil {
		return nil, err
	}

	res := Member{}
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) UpdateMember(ctx context.Context, m *Member) (*Member, error) {
	reqBody, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/network/%s/member/%s", c.BaseURL, m.NetworkID, m.MemberID), reqBody)
	if err != nil {
		return nil, err
	}

	res := Member{}
	if err := c.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) CreateAuthorizedMember(ctx context.Context, networkID, memberID, name string) (*Member, error) {
	m := Member{
		ID:        fmt.Sprintf("%s-%s", networkID, memberID),
		NetworkID: networkID,
		MemberID:  memberID,
		Name:      name,
		Config: MemberConfig{
			Authorized: true,
		},
	}

	reqBody, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/network/%s/member/%s", c.BaseURL, networkID, memberID), reqBody)
	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(ctx, req, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (c *Client) AuthorizeMember(ctx context.Context, networkID, memberID string) (*Member, error) {
	m := Member{
		Config: MemberConfig{
			Authorized: true,
		},
	}

	reqBody, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/network/%s/member/%s", c.BaseURL, networkID, memberID), reqBody)
	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(ctx, req, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (c *Client) DeauthorizeMember(ctx context.Context, networkID, memberID string) (*Member, error) {
	m := Member{
		Config: MemberConfig{
			Authorized: false,
		},
	}

	reqBody, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/network/%s/member/%s", c.BaseURL, networkID, memberID), reqBody)
	if err != nil {
		return nil, err
	}

	if err := c.sendRequest(ctx, req, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (c *Client) DeleteMember(ctx context.Context, m *Member) error {
	req, err := retryablehttp.NewRequest("DELETE", fmt.Sprintf("%s/network/%s/member/%s", c.BaseURL, m.NetworkID, m.MemberID), nil)
	if err != nil {
		return err
	}

	if err := c.sendRequest(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteMemberByID(ctx context.Context, networkID, memberID string) error {
	m := Member{
		NetworkID: networkID,
		MemberID:  memberID,
	}

	return c.DeleteMember(ctx, &m)
}

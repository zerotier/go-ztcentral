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

func (c *Client) GetMembers(ctx context.Context, networkID string) ([]*spec.Member, error) {
	resp, err := c.specClient.GetNetworkMemberList(ctx, networkID)
	if err != nil {
		return nil, err
	}

	var ml []*spec.Member

	return ml, c.decode(resp, &ml)
}

func (c *Client) GetMember(ctx context.Context, networkID, memberID string) (*spec.Member, error) {
	member := &spec.Member{}

	resp, err := c.specClient.GetNetworkMember(ctx, networkID, memberID)
	if err != nil {
		return nil, err
	}

	return member, c.decode(resp, member)
}

func (c *Client) UpdateMember(ctx context.Context, networkID, memberID string, m *spec.Member) (*spec.Member, error) {
	member := &spec.Member{}

	resp, err := c.specClient.UpdateNetworkMember(ctx, networkID, memberID, spec.UpdateNetworkMemberJSONRequestBody(*m))
	if err != nil {
		return nil, err
	}

	return member, c.decode(resp, member)
}

func (c *Client) CreateAuthorizedMember(ctx context.Context, networkID, memberID, name string) (*spec.Member, error) {
	m := &spec.Member{
		NetworkId: &networkID,
		NodeId:    &memberID,
		Name:      &name,
		Config: &spec.MemberConfig{
			Authorized: boolp(true),
		},
	}

	return c.UpdateMember(ctx, networkID, memberID, m)
}

func (c *Client) AuthorizeMember(ctx context.Context, networkID, memberID string) (*spec.Member, error) {
	m := &spec.Member{
		Config: &spec.MemberConfig{
			Authorized: boolp(true),
		},
	}

	return c.UpdateMember(ctx, networkID, memberID, m)
}

func (c *Client) DeauthorizeMember(ctx context.Context, networkID, memberID string) (*spec.Member, error) {
	m := &spec.Member{
		Config: &spec.MemberConfig{
			Authorized: boolp(false),
		},
	}

	return c.UpdateMember(ctx, networkID, memberID, m)
}

func (c *Client) DeleteMember(ctx context.Context, networkID, memberID string) error {
	resp, err := c.specClient.DeleteNetworkMember(ctx, networkID, memberID)
	if err != nil {
		return err
	}

	resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Status was %v: %w", resp.StatusCode, ErrStatus)
	}

	return nil
}

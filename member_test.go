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
// THIS SOFTWARE IS PROVIdED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIdENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package ztcentral

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/zerotier/go-ztcentral/pkg/spec"
	"github.com/zerotier/go-ztcentral/pkg/testutil"
	"github.com/zerotier/go-ztidentity"
)

func TestGetMember(t *testing.T) {
	testutil.NeedsToken(t)

	c, err := NewClient(testutil.InitTokenFromEnv())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	net, err := c.NewNetwork(ctx, "get-member-network", &spec.Network{})
	if err != nil {
		t.Fatal(err)
	}
	defer c.DeleteNetwork(ctx, *net.Config.Id)

	if _, err := c.GetMember(ctx, *net.Config.Id, "123456789"); err == nil {
		t.Fatal("Tried to get a fake member and succeeded")
	}

	aliceID := ztidentity.NewZeroTierIdentity()

	alice, err := c.CreateAuthorizedMember(ctx, *net.Config.Id, aliceID.IDString(), "alice")
	if err != nil {
		t.Fatal(err)
	}

	res, err := c.GetMember(ctx, *net.Config.Id, *alice.NodeId)
	if err != nil {
		t.Fatal(err)
	}

	if *res.NetworkId != *net.Config.Id {
		t.Fatal("network Id of member was not equivalent")
	}

	if *res.NodeId != *alice.NodeId {
		t.Fatal("member Ids were not equivalent")
	}
}

func TestCRUDMembers(t *testing.T) {
	testutil.NeedsToken(t)

	c, err := NewClient(testutil.InitTokenFromEnv())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	net, err := c.NewNetwork(ctx, "get-members-network", &spec.Network{})
	if err != nil {
		t.Fatal(err)
	}
	defer c.DeleteNetwork(ctx, *net.Config.Id)

	users := map[string]ztidentity.ZeroTierIdentity{
		"awk":    ztidentity.NewZeroTierIdentity(),
		"bash":   ztidentity.NewZeroTierIdentity(),
		"cc":     ztidentity.NewZeroTierIdentity(),
		"dpkg":   ztidentity.NewZeroTierIdentity(),
		"edlin":  ztidentity.NewZeroTierIdentity(),
		"finger": ztidentity.NewZeroTierIdentity(),
		"gopher": ztidentity.NewZeroTierIdentity(),
	}

	for name, id := range users {
		_, err := c.CreateAuthorizedMember(ctx, *net.Config.Id, id.IDString(), name)
		if err != nil {
			t.Fatal(err)
		}
	}

	members, err := c.GetMembers(ctx, *net.Config.Id)
	if err != nil {
		t.Fatal(err)
	}

	for _, member := range members {
		id, ok := users[*member.Name]
		if !ok {
			t.Fatal("could not find member in pre-populated table")
		}

		if id.IDString() != *member.NodeId {
			t.Fatalf("Ids were not equal for member %q", *member.Name)
		}
	}

	table := map[string]struct {
		update   func(member *spec.Member)
		validate func(member *spec.Member) error
	}{
		"capabilities": {
			update: func(member *spec.Member) {
				member.Config.Capabilities = &[]int{0, 1, 2}
			},
			validate: func(member *spec.Member) error {
				if !reflect.DeepEqual(*member.Config.Capabilities, []int{0, 1, 2}) {
					return fmt.Errorf("DeepEqual did not succeed on capabilities; was %+v", *member.Config.Capabilities)
				}

				return nil
			},
		},
		"description": {
			update: func(member *spec.Member) {
				member.Description = stringp("updated")
			},
			validate: func(member *spec.Member) error {
				if *member.Description != "updated" {
					return fmt.Errorf("updated value is not present, is: %+v", *member.Description)
				}

				return nil
			},
		},
		"ssoexempt": {
			update: func(member *spec.Member) {
				member.Config.SsoExempt = boolp(true)
			},
			validate: func(member *spec.Member) error {
				if member.Config.SsoExempt == nil || !*member.Config.SsoExempt {
					return fmt.Errorf("Expected SsoExempt to be true, got: %+v", member.Config.SsoExempt)
				}
				return nil
			},
		},
		"ssoexempt_false": {
			update: func(member *spec.Member) {
				member.Config.SsoExempt = boolp(false)
			},
			validate: func(member *spec.Member) error {
				if member.Config.SsoExempt == nil || *member.Config.SsoExempt {
					return fmt.Errorf("Expected SsoExempt to be false, got: %+v", member.Config.SsoExempt)
				}
				return nil
			},
		},
	}

	for _, member := range members {
		for testName, harness := range table {
			harness.update(member)
			updated, err := c.UpdateMember(ctx, *member.NetworkId, *member.NodeId, member)
			if err != nil {
				t.Fatalf("%q: error updating member: %v", testName, err)
			}

			if err := harness.validate(updated); err != nil {
				t.Fatalf("%q: While validating returned object from update call: %v", testName, err)
			}

			newMember, err := c.GetMember(ctx, *member.NetworkId, *member.NodeId)
			if err != nil {
				t.Fatalf("%q: While retrieving updated member from fetch: %v", testName, err)
			}

			if err := harness.validate(newMember); err != nil {
				t.Fatalf("%q: While validating updated member from fetch: %v", testName, err)
			}
		}
	}
}

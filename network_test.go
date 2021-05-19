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
	"testing"
	"time"

	"github.com/zerotier/go-ztcentral/pkg/spec"
	"github.com/zerotier/go-ztcentral/pkg/testutil"
)

func TestNetworkCRUD(t *testing.T) {
	testutil.NeedsToken(t)

	c, err := NewClient(testutil.InitTokenFromEnv())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err = c.GetNetwork(ctx, "8056c2e21c000001")
	if err == nil {
		t.Fatal("Was able to fetch network we don't know about")
	}

	net, err := c.NewNetwork(ctx, "create-network", &spec.Network{})
	if err != nil {
		t.Fatal(err)
	}
	defer c.DeleteNetwork(ctx, *net.Config.Id) // this will fail when the test passes, and that's ok

	res, err := c.GetNetwork(ctx, *net.Config.Id)
	if err != nil {
		t.Fatal("Was able to fetch network we don't know about")
	}

	if *res.Config.Id != *net.Config.Id {
		t.Fatal("Initial returned network configuration was not the same as GetNetwork")
	}

	if *res.Config.Name != *net.Config.Name {
		t.Fatal("Network name was not equal between creation and GetNetwork")
	}

	if err := c.DeleteNetwork(ctx, *net.Config.Id); err != nil {
		t.Fatal(err)
	}

	if _, err := c.GetNetwork(ctx, *net.Config.Id); err == nil {
		t.Fatal("Was able to fetch network we just deleted")
	}
}

func TestNewNetworkWithNetworkConfig(t *testing.T) {
	testutil.NeedsToken(t)

	c, err := NewClient(testutil.InitTokenFromEnv())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// FIXME this should eventually be turned into table tests.
	nc := spec.NetworkConfig{
		Name: stringp("overridden"),
	}

	net, err := c.NewNetwork(ctx, "real", &spec.Network{Config: &nc})
	if err != nil {
		t.Fatal(err)
	}

	if *net.Config.Name != "real" {
		t.Fatal("real name was not overridden during newnetwork")
	}

	getter, err := c.GetNetwork(ctx, *net.Config.Id)
	if err != nil {
		t.Fatal(err)
	}

	if *getter.Config.Name != "real" {
		t.Fatal("real name was not overridden on server side of newnetwork")
	}

	if err := c.DeleteNetwork(ctx, *net.Config.Id); err != nil {
		t.Fatal(err)
	}

	net, err = c.NewNetwork(ctx, "real", &spec.Network{
		RulesSource: stringp("drop;"),
		Config: &spec.NetworkConfig{
			IpAssignmentPools: &[]spec.IPRange{
				{
					IpRangeStart: stringp("10.0.0.2"),
					IpRangeEnd:   stringp("10.0.0.254"),
				},
			},
			Routes: &[]spec.Route{
				{
					Target: stringp("10.0.1.0/24"),
					Via:    stringp("10.0.0.1"),
				},
			},
		}})
	if err != nil {
		t.Fatal(err)
	}

	net, err = c.GetNetwork(ctx, *net.Id)
	if err != nil {
		t.Fatal(err)
	}

	rules, err := c.UpdateNetworkRules(ctx, *net.Id, "drop;")
	if err != nil {
		t.Fatal(err)
	}

	if rules != "drop;" {
		t.Fatal("regurgitated rules were not correct")
	}

	net2, err := c.GetNetwork(ctx, *net.Id)
	if err != nil {
		t.Fatal(err)
	}

	if *net2.RulesSource != "drop;" {
		t.Fatal("rules source was not equal")
	}

	if *net2.Config.Name == "" {
		t.Fatal("name was cleared as a result of rules update")
	}

	if err := c.DeleteNetwork(ctx, *net.Config.Id); err != nil {
		t.Fatal(err)
	}
}

func TestGetNetworks(t *testing.T) {
	testutil.NeedsToken(t)

	c, err := NewClient(testutil.InitTokenFromEnv())

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	networks := map[string]*spec.Network{}

	t.Cleanup(func() {
		for name, net := range networks {
			if err := c.DeleteNetwork(context.Background(), *net.Config.Id); err != nil {
				t.Fatalf("During cleanup of %q: %v", name, err)
			}
		}
	})

	for i := 0; i < 20; i++ {
		name := testutil.RandomString(30, 5)
		net, err := c.NewNetwork(ctx, name, &spec.Network{})
		if err != nil {
			t.Fatal(err)
		}

		networks[name] = net
	}

	res, err := c.GetNetworks(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, network := range res {
		net, ok := networks[*network.Config.Name]
		if !ok {
			continue // not our network, maybe created for some other reason. just ignore
		}

		if *net.Config.Id != *network.Config.Id {
			t.Fatalf("ID mismatch for %q", *network.Config.Name)
		}
	}
}

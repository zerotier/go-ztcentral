// +build integration

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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNetwork(t *testing.T) {
	c := NewClient(os.Getenv("ZEROTIER_CENTRAL_API_KEY"))

	ctx := context.Background()
	res, err := c.GetNetwork(ctx, "8056c2e21c000001")

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, res, "expecting non-nil result")
	assert.Equal(t, "8056c2e21c000001", res.ID, "expecting equal network id")
	assert.Equal(t, "earth.zerotier.net", res.Config.Name, "expecting eqla network name")
}

func TestGetNetworks(t *testing.T) {
	c := NewClient(os.Getenv("ZEROTIER_CENTRAL_API_KEY"))

	ctx := context.Background()
	res, err := c.GetNetworks(ctx)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, res, "expecting non-nil result")
	assert.NotEqual(t, 0, len(res), "expecting non-zero array length")
}

func TestCreateAndDeleteNetwork(t *testing.T) {
	c := NewClient(os.Getenv("ZEROTIER_CENTRAL_API_KEY"))

	ctx := context.Background()
	networkName := "my-test-network"

	res, err := c.NewNetwork(ctx, networkName)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, res, "expecting non-nil result")
	assert.NotEmpty(t, res.ID, "expected network ID to be present")
	assert.Equal(t, networkName, res.Config.Name, "expected equal network names")

	err = c.DeleteNetwork(ctx, res.ID)

	assert.Nil(t, err, "expecting nil error")

	res, err = c.GetNetwork(ctx, res.ID)
	assert.NotNil(t, err, "expecting non-nil error")
	assert.Nil(t, res, "expecting nil result")
}

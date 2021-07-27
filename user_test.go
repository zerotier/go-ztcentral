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

	"github.com/zerotier/go-ztcentral/pkg/testutil"
)

func TestAPITokens(t *testing.T) {
	testutil.NeedsToken(t)

	c, err := NewClient(testutil.InitTokenFromEnv())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	t.Cleanup(cancel)

	token, err := c.RandomToken(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(token) < 32 {
		t.Fatal("token is less than 32 chars")
	}

	user, err := c.User(ctx)
	if err != nil {
		t.Fatalf("Could not fetch API user: %v", err)
	}

	tokenName := testutil.RandomString(30, 30)
	t.Cleanup(func() {
		if err := c.DeleteAPIToken(ctx, *user.Id, tokenName); err != nil {
			t.Fatalf("Could not remove api token: %v", err)
		}
	})

	if err := c.CreateAPIToken(ctx, *user.Id, tokenName, token); err != nil {
		t.Fatalf("While creating API token: %v", err)
	}
}

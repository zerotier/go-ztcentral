package ztcentral

import (
	"context"
	"testing"
	"time"

	"github.com/zerotier/go-ztcentral/pkg/testutil"
)

func TestErrors(t *testing.T) {
	testutil.NeedsToken(t)

	c, err := NewClient(testutil.InitTokenFromEnv())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err = c.GetMember(ctx, "1", "1")
	if err == nil {
		t.Fatal("did not error")
	}
}

func TestUser(t *testing.T) {
	testutil.NeedsToken(t)

	c, err := NewClient(testutil.InitTokenFromEnv())
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	user, err := c.User(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if user == nil || user.Id == nil || len(*user.Id) == 0 {
		t.Fatal("UserID was nil or empty")
	}
}

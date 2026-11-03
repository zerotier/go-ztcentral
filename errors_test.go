package ztcentral

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/zerotier/go-ztcentral/pkg/spec"
)

// newTestClient points a client at handler. BaseURLV1 is read by NewClient,
// so it is swapped for the duration of the call only.
func newTestClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	orig := BaseURLV1
	BaseURLV1 = srv.URL
	defer func() { BaseURLV1 = orig }()

	c, err := NewClient("test-token")
	if err != nil {
		t.Fatal(err)
	}

	return c
}

func respond(status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}
}

func TestAPIErrorFromJSONBody(t *testing.T) {
	c := newTestClient(t, respond(http.StatusLocked, `{"type":"migration_locked","message":"network controller is in migration-lock (read-only) mode"}`))

	_, err := c.UpdateNetwork(context.Background(), "8056c2e21c000001", &spec.Network{})

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}

	if apiErr.StatusCode != http.StatusLocked || apiErr.Type != "migration_locked" {
		t.Fatalf("unexpected APIError: %+v", apiErr)
	}

	if !strings.Contains(err.Error(), "migration-lock") {
		t.Fatalf("server message missing from error: %q", err)
	}

	if !errors.Is(err, ErrStatus) {
		t.Fatal("APIError should still match ErrStatus")
	}
}

func TestAPIErrorFromNonJSONBody(t *testing.T) {
	c := newTestClient(t, respond(http.StatusBadGateway, "<html>bad gateway</html>"))

	_, err := c.GetNetwork(context.Background(), "8056c2e21c000001")

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}

	if string(apiErr.Body) != "<html>bad gateway</html>" {
		t.Fatalf("raw body not kept: %q", apiErr.Body)
	}

	if err.Error() != "Status code 502: status code was not 200" {
		t.Fatalf("unexpected message: %q", err)
	}
}

func TestIsNotFound(t *testing.T) {
	c := newTestClient(t, respond(http.StatusNotFound, `{"type":"not_found","message":"Network not found"}`))
	ctx := context.Background()

	_, err := c.GetMember(ctx, "8056c2e21c000001", "abcdef0123")
	if !IsNotFound(err) {
		t.Fatalf("GetMember: expected not found, got %v", err)
	}

	if err := c.DeleteNetwork(ctx, "8056c2e21c000001"); !IsNotFound(err) {
		t.Fatalf("DeleteNetwork: expected not found, got %v", err)
	}

	if err := c.DeleteMember(ctx, "8056c2e21c000001", "abcdef0123"); !IsNotFound(err) {
		t.Fatalf("DeleteMember: expected not found, got %v", err)
	}

	if IsNotFound(errors.New("Status code 404")) {
		t.Fatal("IsNotFound must only match APIError")
	}
}

func TestDeleteAPITokenReturnsAPIError(t *testing.T) {
	c := newTestClient(t, respond(http.StatusConflict, `{"type":"account_migrated","message":"Make changes there"}`))

	err := c.DeleteAPIToken(context.Background(), "00000000-0000-0000-0000-000000000000", "name")

	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 APIError, got %v", err)
	}
}

func TestTransportErrorDoesNotPanic(t *testing.T) {
	orig := BaseURLV1
	BaseURLV1 = "http://127.0.0.1:1"
	defer func() { BaseURLV1 = orig }()

	c, err := NewClient("test-token")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := c.GetNetwork(context.Background(), "8056c2e21c000001"); err == nil {
		t.Fatal("expected a connection error")
	}
}

func TestConcurrentRequests(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Ratelimit-Limit", "20")
		w.Header().Set("X-Ratelimit-Remaining", "19")
		_, _ = w.Write([]byte(`{"id":"8056c2e21c000001"}`))
	}))

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := c.GetNetwork(context.Background(), "8056c2e21c000001"); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
}

func TestSSOConfigUsesAllowedList(t *testing.T) {
	list := []string{"user@example.com"}
	out, err := json.Marshal(spec.NetworkSSOConfig{AllowedList: &list})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(out), `"allowedList":["user@example.com"]`) {
		t.Fatalf("server expects allowedList, got %s", out)
	}
}

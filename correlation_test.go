// Copyright © 2026 Kluisz.ai, All Rights reserved
// Author: Druhin Abrol <druhin.abrol@kluisz.ai>

package traefik_correlation_id_plugin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	traefik_correlation_id_plugin "github.com/kluisz/traefik-correlation-id-plugin"
)

func TestCorrelation_GeneratesID(t *testing.T) {
	cfg := traefik_correlation_id_plugin.CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := traefik_correlation_id_plugin.New(ctx, next, cfg, "correlation-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeaderNotEmpty(t, req, "X-Klz-Correlation-Id")
}

func TestCorrelation_PreservesExistingID(t *testing.T) {
	cfg := traefik_correlation_id_plugin.CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := traefik_correlation_id_plugin.New(ctx, next, cfg, "correlation-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	existing := "my-existing-correlation-id"
	req.Header.Set("X-Klz-Correlation-Id", existing)

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "X-Klz-Correlation-Id", existing)
}

func TestCorrelation_CustomHeaderName(t *testing.T) {
	cfg := traefik_correlation_id_plugin.CreateConfig()
	cfg.HeaderName = "X-Request-Id"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := traefik_correlation_id_plugin.New(ctx, next, cfg, "correlation-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeaderNotEmpty(t, req, "X-Request-Id")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value for %s: got %q, want %q", key, req.Header.Get(key), expected)
	}
}

func assertHeaderNotEmpty(t *testing.T, req *http.Request, key string) {
	t.Helper()

	if req.Header.Get(key) == "" {
		t.Errorf("expected header %s to be set, but it was empty", key)
	}
}

// Copyright © 2026 Kluisz.ai, All Rights reserved
// Author: Druhin Abrol <druhin.abrol@kluisz.ai>

package traefik_correlation_id_plugin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	traefik_correlation_id_plugin "github.com/navacloud/traefik-correlation-id-plugin"
)

const (
	primaryHeader = "X-Nava-Correlation-Id"
	legacyHeader  = "X-Klz-Correlation-Id"
)

func newHandler(t *testing.T, cfg *traefik_correlation_id_plugin.Config) http.Handler {
	t.Helper()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
	handler, err := traefik_correlation_id_plugin.New(context.Background(), next, cfg, "correlation-plugin")
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func newRequest(t *testing.T) *http.Request {
	t.Helper()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

// TestCorrelation_GeneratesAndDualWrites confirms that when neither header is
// set on the inbound request the plugin generates an ID and writes it to BOTH
// the primary and legacy headers with identical values.
func TestCorrelation_GeneratesAndDualWrites(t *testing.T) {
	handler := newHandler(t, traefik_correlation_id_plugin.CreateConfig())
	req := newRequest(t)

	handler.ServeHTTP(httptest.NewRecorder(), req)

	primary := req.Header.Get(primaryHeader)
	legacy := req.Header.Get(legacyHeader)
	if primary == "" {
		t.Errorf("expected primary header %q to be set", primaryHeader)
	}
	if legacy == "" {
		t.Errorf("expected legacy header %q to be set", legacyHeader)
	}
	if primary != legacy {
		t.Errorf("primary and legacy must match: %q vs %q", primary, legacy)
	}
}

// TestCorrelation_PreservesPrimary confirms an inbound primary header is kept
// and propagated to the legacy header for downstream compatibility.
func TestCorrelation_PreservesPrimary(t *testing.T) {
	handler := newHandler(t, traefik_correlation_id_plugin.CreateConfig())
	req := newRequest(t)
	const existing = "primary-id"
	req.Header.Set(primaryHeader, existing)

	handler.ServeHTTP(httptest.NewRecorder(), req)

	if got := req.Header.Get(primaryHeader); got != existing {
		t.Errorf("primary = %q, want %q", got, existing)
	}
	if got := req.Header.Get(legacyHeader); got != existing {
		t.Errorf("legacy = %q, want %q", got, existing)
	}
}

// TestCorrelation_PromotesLegacy confirms a legacy-only inbound header is
// adopted and mirrored to the primary header (reverse migration support).
func TestCorrelation_PromotesLegacy(t *testing.T) {
	handler := newHandler(t, traefik_correlation_id_plugin.CreateConfig())
	req := newRequest(t)
	const existing = "legacy-id"
	req.Header.Set(legacyHeader, existing)

	handler.ServeHTTP(httptest.NewRecorder(), req)

	if got := req.Header.Get(primaryHeader); got != existing {
		t.Errorf("primary = %q, want %q", got, existing)
	}
	if got := req.Header.Get(legacyHeader); got != existing {
		t.Errorf("legacy = %q, want %q", got, existing)
	}
}

// TestCorrelation_PrimaryWinsOnCollision confirms that when both headers are
// inbound with different values the primary takes precedence and overwrites
// the legacy on the outgoing request.
func TestCorrelation_PrimaryWinsOnCollision(t *testing.T) {
	handler := newHandler(t, traefik_correlation_id_plugin.CreateConfig())
	req := newRequest(t)
	req.Header.Set(primaryHeader, "primary-wins")
	req.Header.Set(legacyHeader, "legacy-loses")

	handler.ServeHTTP(httptest.NewRecorder(), req)

	if got := req.Header.Get(primaryHeader); got != "primary-wins" {
		t.Errorf("primary = %q, want primary-wins", got)
	}
	if got := req.Header.Get(legacyHeader); got != "primary-wins" {
		t.Errorf("legacy = %q, want primary-wins (overwritten)", got)
	}
}

// TestCorrelation_CustomHeaderNames confirms operators can override both names.
func TestCorrelation_CustomHeaderNames(t *testing.T) {
	cfg := traefik_correlation_id_plugin.CreateConfig()
	cfg.HeaderName = "X-Request-Id"
	cfg.LegacyHeaderName = "X-Old-Request-Id"

	handler := newHandler(t, cfg)
	req := newRequest(t)

	handler.ServeHTTP(httptest.NewRecorder(), req)

	primary := req.Header.Get("X-Request-Id")
	legacy := req.Header.Get("X-Old-Request-Id")
	if primary == "" || legacy == "" || primary != legacy {
		t.Errorf("expected matching non-empty headers, got primary=%q legacy=%q", primary, legacy)
	}
}

// Copyright © 2026 Kluisz.ai, All Rights reserved
// Author: Druhin Abrol <druhin.abrol@kluisz.ai>

package traefik_correlation_id_plugin

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// Default header names. The primary header is the canonical post-rebrand name;
// the legacy header is preserved so downstream services not yet migrated still
// see a correlation ID. Both headers are always written with the same value.
const (
	defaultHeaderName       = "X-Nava-Correlation-Id"
	defaultLegacyHeaderName = "X-Klz-Correlation-Id"
)

// Config the plugin configuration.
//
// Backwards-compatibility model: incoming requests may carry the correlation
// ID under either HeaderName or LegacyHeaderName. The plugin reads from the
// primary first, falls back to legacy, and always writes the resolved value
// to BOTH headers on the outgoing request. Set LegacyHeaderName to the empty
// string to disable dual-write once all downstream consumers have migrated.
type Config struct {
	HeaderName       string `yaml:"headerName,omitempty" json:"header_name,omitempty"`
	LegacyHeaderName string `yaml:"legacyHeaderName,omitempty" json:"legacy_header_name,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		HeaderName:       "",
		LegacyHeaderName: "",
	}
}

type Correlation struct {
	next             http.Handler
	name             string
	headerName       string
	legacyHeaderName string
}

// New created a new plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.HeaderName == "" {
		config.HeaderName = defaultHeaderName
	}
	if config.LegacyHeaderName == "" {
		config.LegacyHeaderName = defaultLegacyHeaderName
	}

	return &Correlation{
		next:             next,
		name:             name,
		headerName:       config.HeaderName,
		legacyHeaderName: config.LegacyHeaderName,
	}, nil
}

func (c *Correlation) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// Resolve the correlation ID: prefer the primary header, fall back to the
	// legacy header, generate if neither is set.
	id := request.Header.Get(c.headerName)
	if id == "" && c.legacyHeaderName != "" {
		id = request.Header.Get(c.legacyHeaderName)
	}
	if id == "" {
		b := make([]byte, 16)
		if _, err := rand.Read(b); err == nil {
			id = hex.EncodeToString(b)
		}
	}

	if id != "" {
		request.Header.Set(c.headerName, id)
		if c.legacyHeaderName != "" {
			request.Header.Set(c.legacyHeaderName, id)
		}
	}

	c.next.ServeHTTP(writer, request)
}

// Copyright © 2026 Kluisz.ai, All Rights reserved
// Author: Druhin Abrol <druhin.abrol@kluisz.ai>

package traefik_correlation_id_plugin

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	HeaderName string `yaml:"headerName,omitempty" json:"header_name,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		HeaderName: "",
	}
}

type Correlation struct {
	next       http.Handler
	name       string
	headerName string
}

// New created a new plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.HeaderName == "" {
		config.HeaderName = "X-Klz-Correlation-Id"
	}

	return &Correlation{
		next:       next,
		name:       name,
		headerName: config.HeaderName,
	}, nil
}

func (c *Correlation) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get(c.headerName) == "" {
		b := make([]byte, 16)
		if _, err := rand.Read(b); err == nil {
			request.Header.Set(c.headerName, hex.EncodeToString(b))
		}
	}

	c.next.ServeHTTP(writer, request)
}

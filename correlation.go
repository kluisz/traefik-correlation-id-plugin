// Copyright © 2026 Kluisz.ai, All Rights reserved
// Author: Druhin Abrol <druhin.abrol@kluisz.ai>

package traefik_correlation_id_plugin

import (
	"context"
	"net/http"

	"github.com/kluisz/traefik-correlation-id-plugin/github.com/google/uuid"
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
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
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
	// A Version 7 UUID is a universally unique identifier that is generated using a
	// timestamp, a counter and a cryptographically strong random number. Generally,
	// Version 7 UUIDs have better entropy (i.e. randomness) than Version 1 UUIDs.
	if request.Header.Get(c.headerName) == "" {
		if id, err := uuid.NewV7(); err == nil {
			request.Header.Set(c.headerName, id.String())
		}
	}

	c.next.ServeHTTP(writer, request)
}

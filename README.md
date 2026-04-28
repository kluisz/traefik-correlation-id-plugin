# Traefik Correlation ID Plugin

A Traefik middleware plugin that automatically generates and propagates correlation IDs across requests using UUID v7. Supports dual-header emission for backward compatibility during a header rename.

## Features

- **Automatic UUID v7 Generation** — Generates a unique correlation ID for each request if one doesn't already exist
- **Header Pass-through** — Preserves existing correlation IDs from upstream services
- **Dual-Header Backward Compatibility** — Reads either of two configured header names (preferring `headerName`, falling back to `legacyHeaderName`) and writes the resolved value to both on the outgoing request
- **Configurable Header Names** — Override either or both header names; set `legacyHeaderName: ""` to disable dual-write
- **Zero Dependencies** — Uses only the standard library

## Installation

Add the plugin to your Traefik configuration.

## Configuration

```yaml
http:
  middlewares:
    correlation-id:
      plugin:
        correlation-id:
          headerName: X-Nava-Correlation-Id        # Optional, defaults to X-Nava-Correlation-Id
          legacyHeaderName: X-Klz-Correlation-Id   # Optional, defaults to X-Klz-Correlation-Id; set "" to disable
```

## Usage

Apply the middleware to a route or service:

```yaml
http:
  routers:
    my-router:
      rule: Path(`/api`)
      middlewares:
        - correlation-id
      service: my-service
```

## How It Works

1. The plugin looks for the correlation ID on the inbound request, preferring `headerName` and falling back to `legacyHeaderName`
2. If neither header is present, a new ID is generated from 16 cryptographically random bytes (hex-encoded)
3. The resolved ID is written to both `headerName` and `legacyHeaderName` on the outgoing request, so downstream services that have not yet migrated can still read the legacy header
4. The middleware passes the request to the next handler

## Migration Notes

The defaults reflect the post-rebrand state: `X-Nava-Correlation-Id` is canonical, `X-Klz-Correlation-Id` is the legacy compatibility header. Once all downstream services read the new header, set `legacyHeaderName: ""` to stop emitting the legacy one.

## UUID v7 Format

Version 7 UUIDs use:
- A timestamp component (better sortability)
- A counter for monotonicity
- Cryptographic randomness

This provides better entropy and sorting properties compared to v1 UUIDs.

# Traefik Correlation ID Plugin

A Traefik middleware plugin that automatically generates and propagates correlation IDs across requests using UUID v7.

## Features

- **Automatic UUID v7 Generation** — Generates a unique correlation ID for each request if one doesn't already exist
- **Header Pass-through** — Preserves existing correlation IDs from upstream services
- **Configurable Header Name** — Customize the header name (defaults to `X-Klz-Correlation-Id`)
- **Zero Dependencies** — Uses only the standard library and Google's UUID library

## Installation

Add the plugin to your Traefik configuration.

## Configuration

```yaml
http:
  middlewares:
    correlation-id:
      plugin:
        correlation-id:
          headerName: X-Klz-Correlation-Id  # Optional, defaults to X-Klz-Correlation-Id
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

1. If the request already has a correlation ID header, it is preserved and passed downstream
2. If no correlation ID exists, a new UUID v7 is generated and added to the request
3. The middleware passes the request (with the correlation ID) to the next handler

## UUID v7 Format

Version 7 UUIDs use:
- A timestamp component (better sortability)
- A counter for monotonicity
- Cryptographic randomness

This provides better entropy and sorting properties compared to v1 UUIDs.
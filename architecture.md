# Architecture Overview

This project follows a simple structure to keep business logic decoupled from infrastructure so that proxy, ACME, or database implementations can be swapped easily in the future.

```
VanitySSL/
├── cmd/vanityssl        # Application entrypoint
├── internal/
│   ├── api              # HTTP API for managing customer domains
│   ├── model            # Shared data structures
│   ├── proxy            # Reverse proxy and certificate management
│   └── store            # Storage interfaces and implementations
├── go.mod               # Go module definition
├── README.md            # Project overview
└── architecture.md      # This document
```

## cmd/vanityssl
Contains `main.go`, which reads configuration from environment variables (12‑factor style), constructs the components, and starts the HTTPS proxy. Requests to the hostname specified by `VANITY_API_HOSTNAME` are served by the internal API.

## internal/api
Defines an `API` type that exposes CRUD endpoints for customer domain mappings. An optional bearer token protects the endpoints. The API is served on the main HTTPS port and selected via `VANITY_API_HOSTNAME`.

## internal/model
Holds simple structs shared between packages. The MVP only defines the `Customer` model.

## internal/proxy
Contains the reverse proxy logic. The `Proxy` type uses a `store.Store` to look up which customer owns the domain in the incoming request and injects `X-Customer-ID` and `X-Customer-Domain` headers before forwarding the request to the backend. Certificate handling is abstracted via the `CertManager` interface with a default implementation using `autocert` (`AutoCertManager`).

## internal/store
Defines the `Store` interface for reading and writing domain mappings. The default implementation is `BadgerStore` backed by BadgerDB. A `CachedStore` wraps the backend with an in-memory LRU cache for faster lookups. Because the logic depends only on the interface, other databases (Consul, Etcd, etc.) can be implemented without changing other packages.

## Docker
A `Dockerfile` is provided to build and run the server in a Linux container. Configuration is injected via environment variables at runtime, keeping the image stateless.

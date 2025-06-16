# VanitySSL

VanitySSL is a TLS proxy designed for SaaS platforms that want to offer "vanity" or custom domains with HTTPS support for their customers. It handles certificate management using the ACME protocol and proxies requests to a backend service while injecting headers that identify the customer.

## Overview

A SaaS provider often wants customers to access the service via a custom domain, for example `status.customer.com`. VanitySSL acts as an entry point that:

1. Terminates TLS for each custom domain using SNI to determine which certificate to present.
2. Obtains and renews TLS certificates from an ACME provider such as Let's Encrypt.
3. Proxies requests to the SaaS backend, adding headers so the backend knows which customer the request belongs to. Requests are signed with `X-Vanity-Signature` using the `PROXY_SECRET` environment variable so the backend can verify they came through VanitySSL.
4. Provides an internal API for managing customers and domain mappings.

Customers create a CNAME from their chosen domain to the SaaS endpoint (e.g., `app.saas.com`). VanitySSL sees the incoming SNI and serves the correct certificate for that domain, then forwards the request to the backend service (e.g., `backend.saas.com`).

## Architecture

```
Client --> VanitySSL --> Backend Service
         |           
         |-- ACME (Let's Encrypt) for certificates
         `-- Internal API for managing customers
```

### Components

- **TLS Termination**: Uses Go's TLS stack with SNI support. Certificates are retrieved from storage based on domain and automatically renewed.
- **ACME Client**: The project will use [`golang.org/x/crypto/acme/autocert`](https://pkg.go.dev/golang.org/x/crypto/acme/autocert) for automatic certificate management via ACME. It is part of the Go standard library ecosystem, well-maintained, and widely used for Let's Encrypt integrations.
- **Reverse Proxy**: Requests are forwarded to the backend using [`httputil.ReverseProxy`](https://pkg.go.dev/net/http/httputil#ReverseProxy) from the standard library. It supports modifying requests and responses, making it a good fit for injecting custom headers.
- **Database Interface**: Certificate data and customer domain mappings are stored through an abstract interface so that different key-value stores can be used. The MVP will rely on [BadgerDB](https://github.com/dgraph-io/badger) for local storage. Future implementations may use Consul, Etcd, or a Raft-based store without changing business logic.
- **Internal API**: Exposes CRUD endpoints for customers: `{customerId, domain}`. This allows the SaaS platform to add or remove customer domains at runtime. Configuration (backend address, database directory, etc.) is provided via environment variables.

## Design Decisions

1. **ACME Autocert**: Choosing `autocert` reduces the amount of code needed for certificate management. It handles the ACME challenge/response workflow and certificate renewal out of the box.
2. **Standard Library Reverse Proxy**: `httputil.ReverseProxy` is simple, reliable, and already integrated with Go's HTTP stack. It allows request/response modification, which is enough for adding identifying headers.
3. **Database Abstraction**: Defining a `Store` interface decouples VanitySSL from a specific database. Implementations for BadgerDB and other KV stores can be swapped in as needed.
4. **Minimal Dependencies**: Relying primarily on the Go standard library keeps the code lightweight. External dependencies are chosen carefully for functionality not provided by the standard library (Badger, autocert).

## MVP Implementation Steps

1. **Project Structure**
   - Initialize a Go module.
   - Define the `Store` interface for certificate and domain storage.
   - Implement a BadgerDB-backed store as the default option.

2. **Certificate Management**
   - Integrate `autocert.Manager` to automatically request and renew certificates based on incoming SNI values.
   - Store certificates in the configured `Store` implementation.

3. **Reverse Proxy**
  - Create a proxy handler that looks up the customer via SNI, adds identifying headers (e.g., `X-Customer-ID` and `X-Customer-Domain`), and forwards the request to the configured backend.

4. **Internal API**
   - Implement HTTP endpoints for creating, reading, updating, and deleting customer records.
   - The API should authenticate requests (e.g., via token or IP whitelist) since it modifies domain mappings.

5. **Configuration and Launch**
   - Read environment variables for backend address, database location, ACME email, and other settings.
  - Start the proxy server with HTTPS enabled and serve the internal API on a separate port (port `8081` by default).

6. **Testing and Logging**
   - Add unit tests for the store interface and API logic.
   - Provide logging around certificate renewal and proxy events.

## Future Work

- Implement additional `Store` backends (Consul, Etcd, etc.).
- Support clustering by sharing certificates and domain mappings across nodes via the store.
- Provide metrics (e.g., Prometheus) for monitoring certificate renewals and proxy traffic.
- Expand authentication and security features for the internal API.

VanitySSL's goal is to simplify managing custom domains with TLS so SaaS providers can offer branded endpoints to their customers with minimal operational overhead.


## Running with Docker

Build and run the container:

```sh
docker build -t vanityssl .
docker run -p 80:80 -p 443:443 -p 8081:8081 \
  -e BACKEND_URL=https://backend.internal \
  -e ACME_EMAIL=admin@example.com \
  -e PROXY_SECRET=changeme \
  vanityssl
```

Environment variables configure the backend address, ACME email, optional API token (`API_TOKEN`), database path (`DB_PATH`), and the proxy signing secret (`PROXY_SECRET`). The API is reachable on port `8081`. Port `80` must be reachable for the ACME HTTP-01 challenge. Certificates are stored in the configured database.

A simple test backend is available in `cmd/dummybackend`. It prints request information and verifies the signature.

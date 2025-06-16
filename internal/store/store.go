package store

import "context"

// Store defines methods for storing customer domain mappings and certificates.
// Implementation can use any key-value database.

type Store interface {
	// GetDomain returns the customer ID for the given domain.
	GetDomain(ctx context.Context, domain string) (string, error)
	// SetDomain maps a domain to a customer ID.
	SetDomain(ctx context.Context, domain, customerID string) error
	// DeleteDomain removes a domain mapping.
	DeleteDomain(ctx context.Context, domain string) error
	// ListDomains returns all domain->customer mappings.
	ListDomains(ctx context.Context) (map[string]string, error)

	// GetCert retrieves certificate bytes for the given key. It returns nil
	// if the certificate is not found.
	GetCert(ctx context.Context, key string) ([]byte, error)
	// SetCert stores certificate bytes under the specified key.
	SetCert(ctx context.Context, key string, data []byte) error
	// DeleteCert removes the stored certificate for the key.
	DeleteCert(ctx context.Context, key string) error
}

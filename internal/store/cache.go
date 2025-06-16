package store

import (
	"context"

	lru "github.com/hashicorp/golang-lru/v2"
)

// CachedStore wraps another Store and caches domain lookups.
type CachedStore struct {
	backend Store
	cache   *lru.Cache[string, string]
}

// NewCachedStore creates a CachedStore with the given size.
func NewCachedStore(backend Store, size int) (*CachedStore, error) {
	c, err := lru.New[string, string](size)
	if err != nil {
		return nil, err
	}
	return &CachedStore{backend: backend, cache: c}, nil
}

func (c *CachedStore) GetDomain(ctx context.Context, domain string) (string, error) {
	if val, ok := c.cache.Get(domain); ok {
		return val, nil
	}
	val, err := c.backend.GetDomain(ctx, domain)
	if err != nil {
		return "", err
	}
	if val != "" {
		c.cache.Add(domain, val)
	}
	return val, nil
}

func (c *CachedStore) SetDomain(ctx context.Context, domain, customerID string) error {
	if err := c.backend.SetDomain(ctx, domain, customerID); err != nil {
		return err
	}
	c.cache.Add(domain, customerID)
	return nil
}

func (c *CachedStore) DeleteDomain(ctx context.Context, domain string) error {
	if err := c.backend.DeleteDomain(ctx, domain); err != nil {
		return err
	}
	c.cache.Remove(domain)
	return nil
}

func (c *CachedStore) ListDomains(ctx context.Context) (map[string]string, error) {
	return c.backend.ListDomains(ctx)
}

func (c *CachedStore) GetCert(ctx context.Context, key string) ([]byte, error) {
	return c.backend.GetCert(ctx, key)
}

func (c *CachedStore) SetCert(ctx context.Context, key string, data []byte) error {
	return c.backend.SetCert(ctx, key, data)
}

func (c *CachedStore) DeleteCert(ctx context.Context, key string) error {
	return c.backend.DeleteCert(ctx, key)
}

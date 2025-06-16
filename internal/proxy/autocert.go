package proxy

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/example/vanityssl/internal/store"
	"golang.org/x/crypto/acme/autocert"
)

// AutoCertManager implements CertManager using autocert.Manager.
type AutoCertManager struct {
	Manager *autocert.Manager
}

// storeCache implements autocert.Cache backed by a Store.
type storeCache struct {
	store store.Store
}

func (s *storeCache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := s.store.GetCert(ctx, key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, autocert.ErrCacheMiss
	}
	return data, nil
}

func (s *storeCache) Put(ctx context.Context, key string, data []byte) error {
	return s.store.SetCert(ctx, key, data)
}

func (s *storeCache) Delete(ctx context.Context, key string) error {
	return s.store.DeleteCert(ctx, key)
}

func NewAutoCertManager(st store.Store, email string, hostPolicy autocert.HostPolicy) *AutoCertManager {
	m := &autocert.Manager{
		Cache:      &storeCache{store: st},
		Prompt:     autocert.AcceptTOS,
		Email:      email,
		HostPolicy: hostPolicy,
	}
	return &AutoCertManager{Manager: m}
}

func (a *AutoCertManager) TLSConfig() *tls.Config {
	return &tls.Config{GetCertificate: a.Manager.GetCertificate}
}

// HTTPHandler returns a handler for ACME HTTP-01 challenges.
func (a *AutoCertManager) HTTPHandler(next http.Handler) http.Handler {
	return a.Manager.HTTPHandler(next)
}

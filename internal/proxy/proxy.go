package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/example/vanityssl/internal/store"
)

// Proxy is a reverse proxy that injects customer headers based on SNI.
type Proxy struct {
	Backend     *url.URL
	Store       store.Store
	CertManager CertManager
}

func New(backend *url.URL, st store.Store, cm CertManager) *Proxy {
	return &Proxy{Backend: backend, Store: st, CertManager: cm}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := strings.ToLower(r.Host)
	customerID, err := p.Store.GetDomain(r.Context(), host)
	if err != nil {
		log.Printf("error looking up domain %s: %v", host, err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	r.Header.Set("X-Customer-Domain", host)
	if customerID != "" {
		r.Header.Set("X-Customer-ID", customerID)
	}
	proxy := httputil.NewSingleHostReverseProxy(p.Backend)
	proxy.ErrorLog = log.Default()
	proxy.ServeHTTP(w, r)
}

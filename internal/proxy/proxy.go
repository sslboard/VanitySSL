package proxy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/example/vanityssl/internal/store"
)

// Proxy is a reverse proxy that injects customer headers based on SNI.
type Proxy struct {
	Backend     *url.URL
	Store       store.Store
	CertManager CertManager
	Secret      string
}

func New(backend *url.URL, st store.Store, cm CertManager) *Proxy {
	secret := os.Getenv("PROXY_SECRET")
	return &Proxy{Backend: backend, Store: st, CertManager: cm, Secret: secret}
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
	if p.Secret != "" {
		data := customerID + "\n" + host
		h := hmac.New(sha256.New, []byte(p.Secret))
		h.Write([]byte(data))
		sig := hex.EncodeToString(h.Sum(nil))
		r.Header.Set("X-Vanity-Signature", sig)
	}
	proxy := httputil.NewSingleHostReverseProxy(p.Backend)
	proxy.ErrorLog = log.Default()
	proxy.ServeHTTP(w, r)
}

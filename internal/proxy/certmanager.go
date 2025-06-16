package proxy

import "crypto/tls"

// CertManager abstracts TLS certificate retrieval.
type CertManager interface {
	TLSConfig() *tls.Config
}

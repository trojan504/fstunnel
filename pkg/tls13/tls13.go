// Package tls13 wraps TCP streams in TLS 1.3 only, with AEAD suites
// (AES-GCM or ChaCha20-Poly1305, whichever the CPU prefers).
package tls13

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

// SelfSignedCert makes an ephemeral ECDSA P-256 certificate for one session.
func SelfSignedCert() (tls.Certificate, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}
	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return tls.Certificate{}, err
	}
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: "fstunnel"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, err
	}
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return tls.Certificate{}, err
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	return tls.X509KeyPair(certPEM, keyPEM)
}

// ServerConfig returns a TLS 1.3-only server config.
func ServerConfig(cert tls.Certificate) *tls.Config {
	return &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS13, MaxVersion: tls.VersionTLS13}
}

// ClientConfig returns a TLS 1.3-only client config that skips chain
// verification (the peer is authenticated by the pre-shared tunnel secret
// carried inside the stream, not by a CA).
func ClientConfig() *tls.Config {
	return &tls.Config{MinVersion: tls.VersionTLS13, MaxVersion: tls.VersionTLS13, InsecureSkipVerify: true} //nolint:gosec
}

// Dial opens TCP to addr and completes a TLS 1.3 handshake.
func Dial(addr string, timeout time.Duration) (*tls.Conn, error) {
	d := &net.Dialer{Timeout: timeout}
	raw, err := d.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	c := tls.Client(raw, ClientConfig())
	_ = c.SetDeadline(time.Now().Add(timeout))
	if err := c.Handshake(); err != nil {
		_ = raw.Close()
		return nil, err
	}
	_ = c.SetDeadline(time.Time{})
	return c, nil
}

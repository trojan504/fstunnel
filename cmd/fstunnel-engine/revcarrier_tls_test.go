package main

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"testing"

	"fstunnel/pkg/tls13"
)

func TestRevCarrierTLS13Handshake(t *testing.T) {
	r, err := newRevCarrier(config{role: "iran", secret: "s", port: 1, tls: true})
	if err != nil {
		t.Fatal(err)
	}
	defer r.cancel()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	got := make(chan []byte, 1)
	go func() {
		c, err := ln.Accept()
		if err != nil {
			got <- nil
			return
		}
		sc, err := r.carrierTLSServer(c)
		if err != nil {
			got <- nil
			return
		}
		defer sc.Close()
		b, _ := io.ReadAll(sc)
		got <- b
	}()
	raw, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	tc := tls.Client(raw, tls13.ClientConfig())
	if err := tc.Handshake(); err != nil {
		t.Fatal(err)
	}
	if tc.ConnectionState().Version != tls.VersionTLS13 {
		t.Fatal("not TLS 1.3")
	}
	msg := []byte("carrier-payload")
	if _, err := tc.Write(msg); err != nil {
		t.Fatal(err)
	}
	_ = tc.CloseWrite()
	if b := <-got; !bytes.Equal(b, msg) {
		t.Fatalf("server got %q", b)
	}
}

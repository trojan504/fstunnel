package tls13

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"io"
	"testing"
	"time"
)

func TestTLS13RoundTripAndThroughput(t *testing.T) {
	cert, err := SelfSignedCert()
	if err != nil {
		t.Fatal(err)
	}
	ln, err := tls.Listen("tcp", "127.0.0.1:0", ServerConfig(cert))
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	const size = 64 << 20
	payload := make([]byte, size)
	_, _ = rand.Read(payload)
	done := make(chan []byte, 1)
	go func() {
		c, err := ln.Accept()
		if err != nil {
			done <- nil
			return
		}
		defer c.Close()
		got, _ := io.ReadAll(c)
		done <- got
	}()
	c, err := Dial(ln.Addr().String(), 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if v := c.ConnectionState().Version; v != tls.VersionTLS13 {
		t.Fatalf("negotiated version %x, want TLS 1.3", v)
	}
	start := time.Now()
	if _, err := c.Write(payload); err != nil {
		t.Fatal(err)
	}
	_ = c.CloseWrite()
	got := <-done
	elapsed := time.Since(start).Seconds()
	_ = c.Close()
	if !bytes.Equal(got, payload) {
		t.Fatalf("payload mismatch: got %d bytes", len(got))
	}
	t.Logf("loopback TLS1.3 %.0f Mbps (%d MiB)", float64(size*8)/elapsed/1e6, size>>20)
}

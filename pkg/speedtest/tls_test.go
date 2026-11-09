package speedtest

import (
	"context"
	"testing"
	"time"
)

func TestSpeedOverTLS13AndPlainOnOnePort(t *testing.T) {
	addr, stop, err := StartServer("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer stop()
	ctx := context.Background()
	down, up, err := MeasureSpeed4StreamsTLS(ctx, addr, 200*time.Millisecond, 600*time.Millisecond)
	if err != nil {
		t.Fatalf("tls speed: %v", err)
	}
	if down <= 0 || up <= 0 {
		t.Fatalf("tls speed not measured: down=%.2f up=%.2f", down, up)
	}
	pdown, pup, err := MeasureSpeed4Streams(ctx, addr, 200*time.Millisecond, 600*time.Millisecond)
	if err != nil {
		t.Fatalf("plain speed: %v", err)
	}
	if pdown <= 0 || pup <= 0 {
		t.Fatalf("plain speed not measured: down=%.2f up=%.2f", pdown, pup)
	}
	t.Logf("TLS1.3 down=%.0f up=%.0f Mbps | plain down=%.0f up=%.0f Mbps", down, up, pdown, pup)
}

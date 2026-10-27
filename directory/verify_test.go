package directory

import (
	"strings"
	"testing"
	"time"

	"github.com/xeprotocol/xe/core"
)

func TestVerifyRegistration_ShortNodePeer(t *testing.T) {
	kp, err := core.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair failed: %v", err)
	}

	ts := time.Now().UnixNano()

	// 1. Short NodePeer (< 12 characters) must be rejected
	shortReg := NewRegistration("short-peer", ts, kp)
	if err := VerifyRegistration(shortReg); err == nil {
		t.Fatalf("expected error for short NodePeer (< 12 chars), got nil")
	} else if !strings.Contains(err.Error(), "node_peer too short") {
		t.Fatalf("expected 'node_peer too short' error, got: %v", err)
	}

	// Single character NodePeer
	oneCharReg := NewRegistration("a", ts, kp)
	if err := VerifyRegistration(oneCharReg); err == nil {
		t.Fatalf("expected error for 1-char NodePeer, got nil")
	} else if !strings.Contains(err.Error(), "node_peer too short") {
		t.Fatalf("expected 'node_peer too short' error, got: %v", err)
	}

	// 2. Valid NodePeer (>= 12 characters) must be accepted
	validPeer := "12D3KooWD3g..." // 14 characters
	validReg := NewRegistration(validPeer, ts, kp)
	if err := VerifyRegistration(validReg); err != nil {
		t.Fatalf("expected valid registration to verify cleanly, got: %v", err)
	}
}

package store

import (
	"testing"

	"github.com/xeprotocol/xe/core"
)

func TestMemStore_GetAccountChainLength(t *testing.T) {
	m := NewMemStore()

	// Missing account returns 0, nil
	count, err := m.GetAccountChainLength("missing")
	if err != nil {
		t.Fatalf("unexpected error for missing account: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected count 0, got %d", count)
	}

	// Put chain with 3 blocks
	chain := &core.AccountChain{
		Blocks: []*core.Block{
			{Hash: "b1", Account: "acc1"},
			{Hash: "b2", Account: "acc1"},
			{Hash: "b3", Account: "acc1"},
		},
	}
	if err := m.PutAccountChain("acc1", chain); err != nil {
		t.Fatalf("PutAccountChain failed: %v", err)
	}

	count, err = m.GetAccountChainLength("acc1")
	if err != nil {
		t.Fatalf("GetAccountChainLength failed: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected count 3, got %d", count)
	}
}

func TestBadgerStore_GetAccountChainLength(t *testing.T) {
	dir := t.TempDir()
	s, err := NewBadgerStore(dir)
	if err != nil {
		t.Fatalf("NewBadgerStore failed: %v", err)
	}
	defer s.Close()

	// Missing account returns 0, nil
	count, err := s.GetAccountChainLength("missing")
	if err != nil {
		t.Fatalf("unexpected error for missing account: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected count 0, got %d", count)
	}

	// Put chain with 2 blocks
	chain := &core.AccountChain{
		Blocks: []*core.Block{
			{Hash: "b1", Account: "acc1"},
			{Hash: "b2", Account: "acc1"},
		},
	}
	for _, b := range chain.Blocks {
		if err := s.PutBlock(b); err != nil {
			t.Fatalf("PutBlock failed: %v", err)
		}
	}
	if err := s.PutAccountChain("acc1", chain); err != nil {
		t.Fatalf("PutAccountChain failed: %v", err)
	}

	count, err = s.GetAccountChainLength("acc1")
	if err != nil {
		t.Fatalf("GetAccountChainLength failed: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

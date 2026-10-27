package node

import (
	"testing"
	"time"

	"github.com/xeprotocol/xe/core"
	"github.com/xeprotocol/xe/store"
)

func TestNode_AccountsPagination(t *testing.T) {
	mem := store.NewMemStore()
	l := core.NewLedger(mem, core.LedgerConfig{})

	// Create test accounts with different balances and block counts
	// In addition to genesis account (42T XE), we add:
	// addrB: 500 XE, 2 blocks
	// addrC: 300 XE, 1 block
	// addrA: 300 XE, 1 block (addrA < addrC lexicographically)
	// addrD: 100 XE, 3 blocks
	accounts := []struct {
		addr    string
		balance uint64
		blocks  int
	}{
		{"addrB", 500, 2},
		{"addrC", 300, 1},
		{"addrA", 300, 1},
		{"addrD", 100, 3},
	}

	for _, a := range accounts {
		chain := &core.AccountChain{
			Blocks: make([]*core.Block, a.blocks),
		}
		var lastHash string
		for i := 0; i < a.blocks; i++ {
			h := a.addr + "_b" + string(rune('1'+i))
			b := &core.Block{
				Hash:      h,
				Account:   a.addr,
				Balance:   a.balance,
				Timestamp: time.Now().UnixNano(),
				Asset:     "XE",
			}
			chain.Blocks[i] = b
			if err := mem.PutBlock(b); err != nil {
				t.Fatalf("PutBlock failed: %v", err)
			}
			lastHash = h
		}
		if err := mem.PutAccountChain(a.addr, chain); err != nil {
			t.Fatalf("PutAccountChain failed: %v", err)
		}
		if err := mem.PutFrontier(a.addr, lastHash); err != nil {
			t.Fatalf("PutFrontier failed: %v", err)
		}
		// Register in-memory asset balance
		l.SetAssetBalance(a.addr, "XE", a.balance)
	}

	n := &Node{Ledger: l}

	// 1. Verify GetAllAccounts returns 5 accounts (genesis + 4 test accounts)
	all := n.GetAllAccounts()
	if len(all) != 5 {
		t.Fatalf("expected 5 accounts from GetAllAccounts, got %d", len(all))
	}

	// 2. Genesis account has highest balance, so page at offset 0, limit 1 is genesis
	genPage := n.GetAccountsPage(0, 1)
	if len(genPage) != 1 || genPage[0].Balance < 1000000 {
		t.Fatalf("unexpected genesis page: %+v", genPage)
	}

	// 3. Test GetAccountsPage(1, 2) -> should return addrB (500), then addrA (300, lex before addrC)
	page1 := n.GetAccountsPage(1, 2)
	if len(page1) != 2 {
		t.Fatalf("expected 2 accounts on page 1, got %d", len(page1))
	}
	if page1[0].Address != "addrB" || page1[0].Balance != 500 || page1[0].BlockCount != 2 {
		t.Fatalf("unexpected page1[0]: %+v", page1[0])
	}
	if page1[1].Address != "addrA" || page1[1].Balance != 300 || page1[1].BlockCount != 1 {
		t.Fatalf("unexpected page1[1]: %+v", page1[1])
	}

	// 4. Test GetAccountsPage(3, 2) -> should return addrC (300), then addrD (100)
	page2 := n.GetAccountsPage(3, 2)
	if len(page2) != 2 {
		t.Fatalf("expected 2 accounts on page 2, got %d", len(page2))
	}
	if page2[0].Address != "addrC" || page2[0].Balance != 300 || page2[0].BlockCount != 1 {
		t.Fatalf("unexpected page2[0]: %+v", page2[0])
	}
	if page2[1].Address != "addrD" || page2[1].Balance != 100 || page2[1].BlockCount != 3 {
		t.Fatalf("unexpected page2[1]: %+v", page2[1])
	}

	// 5. Test edge cases
	if empty := n.GetAccountsPage(5, 2); len(empty) != 0 {
		t.Fatalf("expected empty slice when offset >= len, got %d", len(empty))
	}
	if empty := n.GetAccountsPage(10, 5); len(empty) != 0 {
		t.Fatalf("expected empty slice when offset is large, got %d", len(empty))
	}
	if empty := n.GetAccountsPage(0, 0); len(empty) != 0 {
		t.Fatalf("expected empty slice when limit <= 0, got %d", len(empty))
	}
	if empty := n.GetAccountsPage(0, -5); len(empty) != 0 {
		t.Fatalf("expected empty slice when limit < 0, got %d", len(empty))
	}

	// 6. Partial last page (offset 4, limit 5 -> 1 account: addrD)
	last := n.GetAccountsPage(4, 5)
	if len(last) != 1 || last[0].Address != "addrD" {
		t.Fatalf("expected 1 account (addrD) on partial page, got %+v", last)
	}

	// 7. Verify cache consistency: calling getSortedAccountEntries within TTL returns same cache
	entries1 := n.getSortedAccountEntries()
	entries2 := n.getSortedAccountEntries()
	if len(entries1) != len(entries2) {
		t.Fatalf("cached entries length mismatch: %d vs %d", len(entries1), len(entries2))
	}
}

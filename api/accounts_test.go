package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xeprotocol/xe/node"
)

type mockAccountNode struct {
	nodeBackend
	accounts []node.AccountSummary
}

func (m *mockAccountNode) GetAccountsPage(offset, limit int) []node.AccountSummary {
	if limit <= 0 || offset >= len(m.accounts) {
		return []node.AccountSummary{}
	}
	end := offset + limit
	if end > len(m.accounts) {
		end = len(m.accounts)
	}
	return m.accounts[offset:end]
}

func TestHandleGetAllAccounts(t *testing.T) {
	mockAccounts := []node.AccountSummary{
		{Address: "addr1", Balance: 1000, BlockCount: 10, Frontier: "h1"},
		{Address: "addr2", Balance: 500, BlockCount: 5, Frontier: "h2"},
		{Address: "addr3", Balance: 200, BlockCount: 2, Frontier: "h3"},
		{Address: "addr4", Balance: 100, BlockCount: 1, Frontier: "h4"},
	}

	h := &Handler{
		node: &mockAccountNode{accounts: mockAccounts},
	}

	// 1. Test pagination offset=1&limit=2
	req := httptest.NewRequest(http.MethodGet, "/accounts?offset=1&limit=2", nil)
	w := httptest.NewRecorder()
	h.handleGetAllAccounts(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var res []node.AccountSummary
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(res))
	}
	if res[0].Address != "addr2" || res[1].Address != "addr3" {
		t.Fatalf("unexpected accounts returned: %+v", res)
	}

	// 2. Test empty page offset=10
	reqEmpty := httptest.NewRequest(http.MethodGet, "/accounts?offset=10&limit=5", nil)
	wEmpty := httptest.NewRecorder()
	h.handleGetAllAccounts(wEmpty, reqEmpty)

	if wEmpty.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", wEmpty.Code)
	}
	var resEmpty []node.AccountSummary
	if err := json.Unmarshal(wEmpty.Body.Bytes(), &resEmpty); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if len(resEmpty) != 0 {
		t.Fatalf("expected 0 accounts, got %d", len(resEmpty))
	}
}

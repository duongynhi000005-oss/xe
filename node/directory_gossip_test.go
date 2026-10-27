package node

import (
	"testing"
	"time"

	"github.com/xeprotocol/xe/directory"
)

func TestHandleDirectoryGossip_NoPanicOnShortNodePeer(t *testing.T) {
	n := &Node{
		Directory: directory.New(10 * time.Minute),
	}

	regs := make(chan *directory.Registration, 4)

	// 1-character NodePeer that previously panicked with slice bounds out of range [:12]
	regs <- &directory.Registration{
		Account:   "acc_short1",
		NodePeer:  "a",
		Timestamp: time.Now().UnixNano(),
	}

	// 5-character NodePeer
	regs <- &directory.Registration{
		Account:   "acc_short2",
		NodePeer:  "12345",
		Timestamp: time.Now().UnixNano(),
	}

	// Empty NodePeer
	regs <- &directory.Registration{
		Account:   "acc_empty",
		NodePeer:  "",
		Timestamp: time.Now().UnixNano(),
	}

	// Valid length NodePeer
	regs <- &directory.Registration{
		Account:   "acc_valid",
		NodePeer:  "12D3KooWPeerIDValid",
		Timestamp: time.Now().UnixNano(),
	}

	close(regs)

	// Must process without panic
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("handleDirectoryGossip panicked on short NodePeer: %v", r)
		}
	}()

	n.handleDirectoryGossip(regs)
}

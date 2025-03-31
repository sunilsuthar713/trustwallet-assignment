package parser

import (
	"testing"
)

func TestGetCurrentBlock(t *testing.T) {
	p := NewParser()
	block := p.GetCurrentBlock()

	if block <= 0 {
		t.Errorf("Expected a valid block number, got %d", block)
	}
}

func TestSubscribe(t *testing.T) {
	p := NewParser()
	address := "0x123"

	// First subscription should succeed
	if !p.Subscribe(address) {
		t.Errorf("Failed to subscribe address: %s", address)
	}

	// Duplicate subscription should fail
	if p.Subscribe(address) {
		t.Errorf("Duplicate subscription succeeded for address: %s", address)
	}
}

func TestGetTransactions(t *testing.T) {
	p := NewParser()
	address := "0x123"
	p.Subscribe(address)

	transactions := p.GetTransactions(address)
	if transactions == nil {
		t.Errorf("Expected transactions, got nil")
	}
}
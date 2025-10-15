package kernel

import (
	"encoding/hex"
	"testing"
)

func TestTransactionInput(t *testing.T) {
	txBytes, err := hex.DecodeString(coinbaseTxHex)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	tx, err := NewTransaction(txBytes)
	if err != nil {
		t.Fatalf("NewTransaction() error = %v", err)
	}
	defer tx.Destroy()

	// Get transaction input
	inputView, err := tx.GetInput(0)
	if err != nil {
		t.Fatalf("GetInput(0) error = %v", err)
	}
	if inputView == nil {
		t.Fatal("GetInput(0) returned nil")
	}

	// Test Copy()
	copiedInput := inputView.Copy()
	if copiedInput == nil {
		t.Fatal("Copy() returned nil")
	}
	defer copiedInput.Destroy()

	// Test GetOutPoint()
	outPoint := inputView.GetOutPoint()
	if outPoint == nil {
		t.Fatal("GetOutPoint() returned nil")
	}

	// Verify out point from copied input matches
	copiedOutPoint := copiedInput.GetOutPoint()
	if copiedOutPoint == nil {
		t.Fatal("GetOutPoint() on copied input returned nil")
	}

	// Compare txid bytes from both out points
	originalTxid := outPoint.GetTxid()
	copiedTxid := copiedOutPoint.GetTxid()

	if originalTxid.Bytes() != copiedTxid.Bytes() {
		t.Errorf("OutPoint txids differ: %x != %x", originalTxid.Bytes(), copiedTxid.Bytes())
	}

	// Compare index from both out points
	if outPoint.GetIndex() != copiedOutPoint.GetIndex() {
		t.Errorf("OutPoint indices differ: %d != %d", outPoint.GetIndex(), copiedOutPoint.GetIndex())
	}
}

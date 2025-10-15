package kernel

import (
	"encoding/hex"
	"testing"
)

func TestTransactionOutPoint(t *testing.T) {
	txBytes, err := hex.DecodeString(coinbaseTxHex)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	tx, err := NewTransaction(txBytes)
	if err != nil {
		t.Fatalf("NewTransaction() error = %v", err)
	}
	defer tx.Destroy()

	// Get transaction input and its out point
	inputView, err := tx.GetInput(0)
	if err != nil {
		t.Fatalf("GetInput(0) error = %v", err)
	}

	outPointView := inputView.GetOutPoint()
	if outPointView == nil {
		t.Fatal("GetOutPoint() returned nil")
	}

	// Test Copy()
	copiedOutPoint := outPointView.Copy()
	if copiedOutPoint == nil {
		t.Fatal("Copy() returned nil")
	}
	defer copiedOutPoint.Destroy()

	// Test GetTxid()
	originalTxid := outPointView.GetTxid()
	if originalTxid == nil {
		t.Fatal("GetTxid() returned nil")
	}

	copiedTxid := copiedOutPoint.GetTxid()
	if copiedTxid == nil {
		t.Fatal("GetTxid() on copied out point returned nil")
	}

	if originalTxid.Bytes() != copiedTxid.Bytes() {
		t.Errorf("Txid bytes differ: %x != %x", originalTxid.Bytes(), copiedTxid.Bytes())
	}

	// Test GetIndex()
	originalIndex := outPointView.GetIndex()
	copiedIndex := copiedOutPoint.GetIndex()

	if originalIndex != copiedIndex {
		t.Errorf("Indices differ: %d != %d", originalIndex, copiedIndex)
	}

	// Verify the index is what we expect for a coinbase transaction (0xffffffff)
	if originalIndex != 0xffffffff {
		t.Errorf("Expected coinbase index 0xffffffff, got %d", originalIndex)
	}
}

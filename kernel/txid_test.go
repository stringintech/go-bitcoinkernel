package kernel

import (
	"encoding/hex"
	"testing"
)

func TestTxid(t *testing.T) {
	txBytes, err := hex.DecodeString(coinbaseTxHex)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	tx, err := NewTransaction(txBytes)
	if err != nil {
		t.Fatalf("NewTransaction() error = %v", err)
	}
	defer tx.Destroy()
	txid := tx.GetTxid()
	if txid == nil {
		t.Fatal("GetTxid() returned nil")
	}

	// Test Bytes()
	txidBytes := txid.Bytes()
	if txidBytes == [32]byte{} {
		t.Error("Txid.Bytes() returned empty bytes")
	}

	// Test Copy()
	copiedTxid := txid.Copy()
	defer copiedTxid.Destroy()

	if txid.Bytes() != copiedTxid.Bytes() {
		t.Errorf("Copied txid bytes differ: %x != %x", txid.Bytes(), copiedTxid.Bytes())
	}

	// Test Equals()
	if !txid.Equals(copiedTxid) {
		t.Error("txid.Equals(copiedTxid) = false, want true")
	}
}

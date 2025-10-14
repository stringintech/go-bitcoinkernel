package kernel

import (
	"encoding/hex"
	"testing"
)

func TestTxid(t *testing.T) {
	txHex := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff08044c86041b020602ffffffff0100f2052a010000004341041b0e8c2567c12536aa13357b79a073dc4444acb83c4ec7a0e2f99dd7457516c5817242da796924ca4e99947d087fedf9ce467cb9f7c6287078f801df276fdf84ac00000000"
	txBytes, err := hex.DecodeString(txHex)
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

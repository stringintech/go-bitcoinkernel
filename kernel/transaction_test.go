package kernel

import (
	"encoding/hex"
	"errors"
	"testing"
)

func TestInvalidTransactionData(t *testing.T) {
	// Test with empty data
	_, err := NewTransactionFromRaw([]byte{})
	if !errors.Is(err, ErrInvalidTransactionData) {
		t.Errorf("Expected ErrInvalidTransactionData, got %v", err)
	}

	// Test with invalid data
	_, err = NewTransactionFromRaw([]byte{0x00, 0x01, 0x02})
	if !errors.Is(err, ErrTransactionCreation) {
		t.Errorf("Expected ErrTransactionCreation, got %v", err)
	}
}

func TestTransactionFromRaw(t *testing.T) {
	txHex := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff08044c86041b020602ffffffff0100f2052a010000004341041b0e8c2567c12536aa13357b79a073dc4444acb83c4ec7a0e2f99dd7457516c5817242da796924ca4e99947d087fedf9ce467cb9f7c6287078f801df276fdf84ac00000000"
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	tx, err := NewTransactionFromRaw(txBytes)
	if err != nil {
		t.Fatalf("NewTransactionFromRaw() error = %v", err)
	}
	if tx == nil {
		t.Fatal("Transaction is nil")
	}
	defer tx.Destroy()

	if tx.ptr == nil {
		t.Error("Transaction pointer is nil")
	}
}

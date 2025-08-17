package kernel

import (
	"encoding/hex"
	"errors"
	"testing"
)

func TestInvalidTransactionData(t *testing.T) {
	// Test with empty data
	_, err := NewTransactionFromRaw([]byte{})
	if !errors.Is(err, ErrEmptyTransactionData) {
		t.Errorf("Expected ErrEmptyTransactionData, got %v", err)
	}

	// Test with invalid data
	_, err = NewTransactionFromRaw([]byte{0x00, 0x01, 0x02})
	if !errors.Is(err, ErrKernelTransactionCreate) {
		t.Errorf("Expected ErrKernelTransactionCreate, got %v", err)
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

func TestTransactionCopy(t *testing.T) {
	txHex := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff08044c86041b020602ffffffff0100f2052a010000004341041b0e8c2567c12536aa13357b79a073dc4444acb83c4ec7a0e2f99dd7457516c5817242da796924ca4e99947d087fedf9ce467cb9f7c6287078f801df276fdf84ac00000000"
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	tx, err := NewTransactionFromRaw(txBytes)
	if err != nil {
		t.Fatalf("NewTransactionFromRaw() error = %v", err)
	}
	defer tx.Destroy()

	// Test copying transaction
	txCopy, err := tx.Copy()
	if err != nil {
		t.Fatalf("Transaction.Copy() error = %v", err)
	}
	if txCopy == nil {
		t.Fatal("Copied transaction is nil")
	}
	defer txCopy.Destroy()

	if txCopy.ptr == nil {
		t.Error("Copied transaction pointer is nil")
	}
}

func TestTransactionCountInputsOutputs(t *testing.T) {
	txHex := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff08044c86041b020602ffffffff0100f2052a010000004341041b0e8c2567c12536aa13357b79a073dc4444acb83c4ec7a0e2f99dd7457516c5817242da796924ca4e99947d087fedf9ce467cb9f7c6287078f801df276fdf84ac00000000"
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	tx, err := NewTransactionFromRaw(txBytes)
	if err != nil {
		t.Fatalf("NewTransactionFromRaw() error = %v", err)
	}
	defer tx.Destroy()

	// Test counting inputs (this is a coinbase transaction with 1 input)
	inputCount, err := tx.CountInputs()
	if err != nil {
		t.Fatalf("Transaction.CountInputs() error = %v", err)
	}
	if inputCount != 1 {
		t.Errorf("Expected 1 input, got %d", inputCount)
	}

	// Test counting outputs (this transaction has 1 output)
	outputCount, err := tx.CountOutputs()
	if err != nil {
		t.Fatalf("Transaction.CountOutputs() error = %v", err)
	}
	if outputCount != 1 {
		t.Errorf("Expected 1 output, got %d", outputCount)
	}
}

func TestTransactionGetOutputAt(t *testing.T) {
	txHex := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff08044c86041b020602ffffffff0100f2052a010000004341041b0e8c2567c12536aa13357b79a073dc4444acb83c4ec7a0e2f99dd7457516c5817242da796924ca4e99947d087fedf9ce467cb9f7c6287078f801df276fdf84ac00000000"
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	tx, err := NewTransactionFromRaw(txBytes)
	if err != nil {
		t.Fatalf("NewTransactionFromRaw() error = %v", err)
	}
	defer tx.Destroy()

	// Test getting output at index 0
	output, err := tx.GetOutputAt(0)
	if err != nil {
		t.Fatalf("Transaction.GetOutputAt(0) error = %v", err)
	}
	if output == nil {
		t.Fatal("Output is nil")
	}
	defer output.Destroy()

	if output.ptr == nil {
		t.Error("Output pointer is nil")
	}
}

func TestTransactionToBytes(t *testing.T) {
	txHex := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff08044c86041b020602ffffffff0100f2052a010000004341041b0e8c2567c12536aa13357b79a073dc4444acb83c4ec7a0e2f99dd7457516c5817242da796924ca4e99947d087fedf9ce467cb9f7c6287078f801df276fdf84ac00000000"
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	tx, err := NewTransactionFromRaw(txBytes)
	if err != nil {
		t.Fatalf("NewTransactionFromRaw() error = %v", err)
	}
	defer tx.Destroy()

	// Test serializing transaction back to bytes
	serialized, err := tx.ToBytes()
	if err != nil {
		t.Fatalf("Transaction.ToBytes() error = %v", err)
	}

	if len(serialized) == 0 {
		t.Error("Serialized transaction is empty")
	}

	// The serialized bytes should match the original
	if hex.EncodeToString(serialized) != txHex {
		t.Errorf("Serialized transaction doesn't match original.\nExpected: %s\nGot: %s", txHex, hex.EncodeToString(serialized))
	}
}

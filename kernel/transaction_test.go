package kernel

import (
	"encoding/hex"
	"errors"
	"testing"
)

// coinbaseTxHex is a serialized coinbase transaction for testing
const coinbaseTxHex = "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff08044c86041b020602ffffffff0100f2052a010000004341041b0e8c2567c12536aa13357b79a073dc4444acb83c4ec7a0e2f99dd7457516c5817242da796924ca4e99947d087fedf9ce467cb9f7c6287078f801df276fdf84ac00000000"

func TestInvalidTransactionData(t *testing.T) {
	// Test with invalid data
	_, err := NewTransaction([]byte{0x00, 0x01, 0x02})
	var internalErr *InternalError
	if !errors.As(err, &internalErr) {
		t.Errorf("Expected InternalError, got %v", err)
	}
}

func TestTransaction(t *testing.T) {
	txBytes, err := hex.DecodeString(coinbaseTxHex)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	tx, err := NewTransaction(txBytes)
	if err != nil {
		t.Fatalf("NewTransaction() error = %v", err)
	}
	if tx == nil {
		t.Fatal("Transaction is nil")
	}
	defer tx.Destroy()

	if tx.handle.ptr == nil {
		t.Error("Transaction pointer is nil")
	}

	// Test Copy()
	txCopy := tx.Copy()
	if txCopy == nil {
		t.Fatal("Copied transaction is nil")
	}
	defer txCopy.Destroy()

	if txCopy.handle.ptr == nil {
		t.Error("Copied transaction pointer is nil")
	}

	// Test CountInputs() (this is a coinbase transaction with 1 input)
	inputCount := tx.CountInputs()
	if inputCount != 1 {
		t.Errorf("Expected 1 input, got %d", inputCount)
	}

	// Test CountOutputs() (this transaction has 1 output)
	outputCount := tx.CountOutputs()
	if outputCount != 1 {
		t.Errorf("Expected 1 output, got %d", outputCount)
	}

	// Test GetInput()
	input, err := tx.GetInput(0)
	if err != nil {
		t.Fatalf("GetInput(0) error = %v", err)
	}
	if input == nil {
		t.Error("Input is nil")
	}
	_, err = tx.GetInput(inputCount)
	if !errors.Is(err, ErrKernelIndexOutOfBounds) {
		t.Errorf("Expected ErrKernelIndexOutOfBounds for out of bounds input, got %v", err)
	}

	// Test GetOutput()
	output, err := tx.GetOutput(0)
	if err != nil {
		t.Fatalf("GetOutput(0) error = %v", err)
	}
	if output == nil {
		t.Fatal("Output is nil")
	}
	_, err = tx.GetOutput(outputCount)
	if !errors.Is(err, ErrKernelIndexOutOfBounds) {
		t.Errorf("Expected ErrKernelIndexOutOfBounds for out of bounds output, got %v", err)
	}

	// Test GetTxid()
	_ = tx.GetTxid()

	// Test Bytes()
	serialized, err := tx.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error = %v", err)
	}

	if len(serialized) == 0 {
		t.Error("Serialized transaction is empty")
	}

	// The serialized bytes should match the original
	if hex.EncodeToString(serialized) != coinbaseTxHex {
		t.Errorf("Serialized transaction doesn't match original.\nExpected: %s\nGot: %s", coinbaseTxHex, hex.EncodeToString(serialized))
	}
}

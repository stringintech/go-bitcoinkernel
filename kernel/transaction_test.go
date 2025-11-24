package kernel

import (
	"encoding/hex"
	"errors"
	"slices"
	"testing"
)

// coinbaseTxHex is a serialized coinbase transaction for testing
const coinbaseTxHex = "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff08044c86041b020602ffffffff0100f2052a010000004341041b0e8c2567c12536aa13357b79a073dc4444acb83c4ec7a0e2f99dd7457516c5817242da796924ca4e99947d087fedf9ce467cb9f7c6287078f801df276fdf84ac00000000"

func TestInvalidTransactionData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"invalid bytes", []byte{0x00, 0x01, 0x02}},
		{"nil slice", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTransaction(tt.data)
			var internalErr *InternalError
			if !errors.As(err, &internalErr) {
				t.Errorf("Expected InternalError, got %v", err)
			}
		})
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

	t.Run("Copy", func(t *testing.T) {
		txCopy := tx.Copy()
		if txCopy == nil {
			t.Fatal("Copied transaction is nil")
		}
		defer txCopy.Destroy()

		if txCopy.handle.ptr == nil {
			t.Error("Copied transaction pointer is nil")
		}
		if txCopy.handle.ptr == tx.handle.ptr {
			t.Error("Copied transaction pointer should be different from original")
		}
	})

	t.Run("GetTxid", func(t *testing.T) {
		txid := tx.GetTxid()
		if txid == nil {
			t.Error("Txid is nil")
		}
	})

	t.Run("Bytes", func(t *testing.T) {
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
	})

	t.Run("CountInputs", func(t *testing.T) {
		// This is a coinbase transaction with 1 input
		inputCount := tx.CountInputs()
		if inputCount != 1 {
			t.Errorf("Expected 1 input, got %d", inputCount)
		}
	})

	t.Run("GetInput", func(t *testing.T) {
		input, err := tx.GetInput(0)
		if err != nil {
			t.Fatalf("GetInput(0) error = %v", err)
		}
		if input == nil {
			t.Error("Input is nil")
		}

		_, err = tx.GetInput(tx.CountInputs())
		if !errors.Is(err, ErrKernelIndexOutOfBounds) {
			t.Errorf("Expected ErrKernelIndexOutOfBounds for out of bounds input, got %v", err)
		}
	})

	t.Run("Inputs", func(t *testing.T) {
		count := len(slices.Collect(tx.Inputs()))
		if count != 1 {
			t.Errorf("Expected to iterate over 1 input, got %d", count)
		}
	})

	t.Run("InputsRange", func(t *testing.T) {
		count := len(slices.Collect(tx.InputsRange(0, 1000)))
		if count != 1 {
			t.Errorf("Expected to iterate over 1 input, got %d", count)
		}

		count = len(slices.Collect(tx.InputsRange(1, 2)))
		if count != 0 {
			t.Errorf("Expected to iterate over 0 inputs, got %d", count)
		}
	})

	t.Run("InputsFrom", func(t *testing.T) {
		count := len(slices.Collect(tx.InputsFrom(0)))
		if count != 1 {
			t.Errorf("Expected to iterate over 1 input, got %d", count)
		}

		count = len(slices.Collect(tx.InputsFrom(1)))
		if count != 0 {
			t.Errorf("Expected to iterate over 0 inputs, got %d", count)
		}
	})

	t.Run("CountOutputs", func(t *testing.T) {
		// This transaction has 1 output
		outputCount := tx.CountOutputs()
		if outputCount != 1 {
			t.Errorf("Expected 1 output, got %d", outputCount)
		}
	})

	t.Run("GetOutput", func(t *testing.T) {
		output, err := tx.GetOutput(0)
		if err != nil {
			t.Fatalf("GetOutput(0) error = %v", err)
		}
		if output == nil {
			t.Fatal("Output is nil")
		}

		_, err = tx.GetOutput(tx.CountOutputs())
		if !errors.Is(err, ErrKernelIndexOutOfBounds) {
			t.Errorf("Expected ErrKernelIndexOutOfBounds for out of bounds output, got %v", err)
		}
	})

	t.Run("Outputs", func(t *testing.T) {
		count := len(slices.Collect(tx.Outputs()))
		if count != 1 {
			t.Errorf("Expected to iterate over 1 output, got %d", count)
		}
	})

	t.Run("OutputsRange", func(t *testing.T) {
		count := len(slices.Collect(tx.OutputsRange(0, 1000)))
		if count != 1 {
			t.Errorf("Expected to iterate over 1 output, got %d", count)
		}

		count = len(slices.Collect(tx.OutputsRange(1, 2)))
		if count != 0 {
			t.Errorf("Expected to iterate over 0 outputs, got %d", count)
		}
	})

	t.Run("OutputsFrom", func(t *testing.T) {
		count := len(slices.Collect(tx.OutputsFrom(0)))
		if count != 1 {
			t.Errorf("Expected to iterate over 1 output, got %d", count)
		}

		count = len(slices.Collect(tx.OutputsFrom(1)))
		if count != 0 {
			t.Errorf("Expected to iterate over 0 outputs, got %d", count)
		}
	})
}

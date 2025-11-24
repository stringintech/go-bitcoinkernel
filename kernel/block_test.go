package kernel

import (
	"encoding/hex"
	"errors"
	"slices"
	"testing"
)

func TestInvalidBlockData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"invalid bytes", []byte{0x00, 0x01, 0x02}},
		{"nil slice", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewBlock(tt.data)
			var internalErr *InternalError
			if !errors.As(err, &internalErr) {
				t.Errorf("Expected InternalError, got %v", err)
			}
		})
	}
}

func TestBlock(t *testing.T) {
	// Complete Bitcoin mainnet genesis block (285 bytes)
	genesisHex := "0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c0101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000"
	genesisBytes, err := hex.DecodeString(genesisHex)
	if err != nil {
		t.Fatalf("Failed to decode genesis hex: %v", err)
	}

	block, err := NewBlock(genesisBytes)
	if err != nil {
		t.Fatalf("NewBlock() error = %v", err)
	}
	defer block.Destroy()

	t.Run("Hash", func(t *testing.T) {
		hash := block.Hash()
		defer hash.Destroy()

		// Expected genesis block hash (reversed byte order for display)
		expectedHash := "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"
		hashBytes := hash.Bytes()
		actualHashHex := hex.EncodeToString(reverseBytes(hashBytes[:]))
		if actualHashHex != expectedHash {
			t.Errorf("Expected hash %s, got %s", expectedHash, actualHashHex)
		}
	})

	t.Run("Bytes", func(t *testing.T) {
		data, err := block.Bytes()
		if err != nil {
			t.Fatalf("Block.Bytes() error = %v", err)
		}

		if len(data) != len(genesisBytes) {
			t.Errorf("Expected data length %d, got %d", len(genesisBytes), len(data))
		}

		hexStr := hex.EncodeToString(data)
		if hexStr != genesisHex {
			t.Errorf("Expected data hex %s, got %s", genesisHex, hexStr)
		}
	})

	t.Run("Copy", func(t *testing.T) {
		blockCopy := block.Copy()
		if blockCopy == nil {
			t.Fatal("Copied block is nil")
		}
		defer blockCopy.Destroy()

		if blockCopy.ptr == nil {
			t.Error("Copied block pointer is nil")
		}
	})

	t.Run("CountTransactions", func(t *testing.T) {
		// Genesis block has 1 transaction
		txCount := block.CountTransactions()
		if txCount != 1 {
			t.Errorf("Expected 1 transaction, got %d", txCount)
		}
	})

	t.Run("GetTransactionAt", func(t *testing.T) {
		tx, err := block.GetTransactionAt(0)
		if err != nil {
			t.Fatalf("Block.GetTransactionAt(0) error = %v", err)
		}
		if tx == nil {
			t.Fatal("Transaction is nil")
		}
		if tx.ptr == nil {
			t.Error("Transaction pointer is nil")
		}
	})

	t.Run("Transactions", func(t *testing.T) {
		count := len(slices.Collect(block.Transactions()))
		if count != 1 {
			t.Errorf("Expected to iterate over 1 transaction, got %d", count)
		}
	})

	t.Run("TransactionsRange", func(t *testing.T) {
		count := len(slices.Collect(block.TransactionsRange(0, 1000)))
		if count != 1 {
			t.Errorf("Expected to iterate over 1 transaction, got %d", count)
		}

		count = len(slices.Collect(block.TransactionsRange(1, 2)))
		if count != 0 {
			t.Errorf("Expected to iterate over 0 transactions, got %d", count)
		}
	})

	t.Run("TransactionsFrom", func(t *testing.T) {
		count := len(slices.Collect(block.TransactionsFrom(0)))
		if count != 1 {
			t.Errorf("Expected to iterate over 1 transaction, got %d", count)
		}

		count = len(slices.Collect(block.TransactionsFrom(1)))
		if count != 0 {
			t.Errorf("Expected to iterate over 0 transactions, got %d", count)
		}
	})
}

func reverseBytes(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[len(data)-1-i] = b
	}
	return result
}

package kernel

import (
	"encoding/hex"
	"errors"
	"testing"
)

func TestInvalidBlockData(t *testing.T) {
	_, err := NewBlock([]byte{0x00, 0x01, 0x02})
	var internalErr *InternalError
	if !errors.As(err, &internalErr) {
		t.Errorf("Expected InternalError, got %v", err)
	}
}

func TestBlockFromRaw(t *testing.T) {
	// Complete Bitcoin mainnet genesis block (285 bytes)
	genesisHex := "0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c0101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000"
	genesisBytes, err := hex.DecodeString(genesisHex)
	if err != nil {
		t.Fatalf("Failed to decode genesis hex: %v", err)
	}

	block, err := NewBlock(genesisBytes)
	if err != nil {
		t.Fatalf("NewBlockFromRaw() error = %v", err)
	}
	defer block.Destroy()

	// Test getting block hash
	hash := block.Hash()
	defer hash.Destroy()

	// Expected genesis block hash (reversed byte order for display)
	expectedHash := "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"
	hashBytes := hash.Bytes()
	actualHashHex := hex.EncodeToString(reverseBytes(hashBytes[:]))
	if actualHashHex != expectedHash {
		t.Logf("Actual hash: %s", actualHashHex)
		t.Logf("Expected hash: %s", expectedHash)
	}

	// Test getting the serialized block
	data, err := block.Bytes()
	if err != nil {
		t.Fatalf("Block.Bytes() error = %v", err)
	}

	if len(data) != len(genesisBytes) {
		t.Errorf("Expected data length %d, got %d", len(genesisBytes), len(data))
	}

	hexStr := hex.EncodeToString(data)
	if hexStr != genesisHex {
		t.Logf("Expected data hex: %s, got %s", genesisHex, hexStr)
	}
}

func TestBlockCopy(t *testing.T) {
	genesisHex := "0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c0101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000"
	genesisBytes, err := hex.DecodeString(genesisHex)
	if err != nil {
		t.Fatalf("Failed to decode genesis hex: %v", err)
	}

	block, err := NewBlock(genesisBytes)
	if err != nil {
		t.Fatalf("NewBlockFromRaw() error = %v", err)
	}
	defer block.Destroy()

	// Test copying block
	blockCopy := block.Copy()
	if blockCopy == nil {
		t.Fatal("Copied block is nil")
	}
	defer blockCopy.Destroy()

	if blockCopy.ptr == nil {
		t.Error("Copied block pointer is nil")
	}
}

func TestBlockCountTransactions(t *testing.T) {
	genesisHex := "0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c0101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000"
	genesisBytes, err := hex.DecodeString(genesisHex)
	if err != nil {
		t.Fatalf("Failed to decode genesis hex: %v", err)
	}

	block, err := NewBlock(genesisBytes)
	if err != nil {
		t.Fatalf("NewBlockFromRaw() error = %v", err)
	}
	defer block.Destroy()

	// Test counting transactions (genesis block has 1 transaction)
	txCount := block.CountTransactions()
	if txCount != 1 {
		t.Errorf("Expected 1 transaction, got %d", txCount)
	}
}

func TestBlockGetTransactionAt(t *testing.T) {
	genesisHex := "0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c0101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000"
	genesisBytes, err := hex.DecodeString(genesisHex)
	if err != nil {
		t.Fatalf("Failed to decode genesis hex: %v", err)
	}

	block, err := NewBlock(genesisBytes)
	if err != nil {
		t.Fatalf("NewBlockFromRaw() error = %v", err)
	}
	defer block.Destroy()

	// Test getting transaction at index 0
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
}

func reverseBytes(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[len(data)-1-i] = b
	}
	return result
}

package kernel

import (
	"encoding/hex"
	"errors"
	"testing"
)

// genesisHeaderHex is the Bitcoin mainnet genesis block header (80 bytes)
const genesisHeaderHex = "0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c"

func TestInvalidBlockHeaderData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"invalid bytes", []byte{0x00, 0x01, 0x02}},
		{"nil slice", nil},
		{"too short", make([]byte, 79)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewBlockHeader(tt.data)
			var internalErr *InternalError
			if !errors.As(err, &internalErr) {
				t.Errorf("Expected InternalError, got %v", err)
			}
		})
	}
}

func TestBlockHeader(t *testing.T) {
	headerBytes, err := hex.DecodeString(genesisHeaderHex)
	if err != nil {
		t.Fatalf("Failed to decode genesis header hex: %v", err)
	}

	header, err := NewBlockHeader(headerBytes)
	if err != nil {
		t.Fatalf("NewBlockHeader() error = %v", err)
	}
	if header == nil {
		t.Fatal("BlockHeader is nil")
	}
	defer header.Destroy()

	if header.ptr == nil {
		t.Error("BlockHeader pointer is nil")
	}

	t.Run("Copy", func(t *testing.T) {
		headerCopy := header.Copy()
		if headerCopy == nil {
			t.Fatal("Copied header is nil")
		}
		defer headerCopy.Destroy()

		if headerCopy.ptr == nil {
			t.Error("Copied header pointer is nil")
		}
		if headerCopy.ptr == header.ptr {
			t.Error("Copied header pointer should be different from original")
		}
	})

	t.Run("Hash", func(t *testing.T) {
		hash := header.Hash()
		if hash == nil {
			t.Fatal("Hash is nil")
		}
		defer hash.Destroy()

		// Expected genesis block hash (reversed byte order for display)
		expectedHash := "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"
		actualHashHex := hash.String()
		if actualHashHex != expectedHash {
			t.Errorf("Expected hash %s, got %s", expectedHash, actualHashHex)
		}
	})

	t.Run("PrevHash", func(t *testing.T) {
		prevHash := header.PrevHash()
		prevHashStr := prevHash.String()

		// Genesis block has all-zero previous hash
		if prevHashStr != "0000000000000000000000000000000000000000000000000000000000000000" {
			t.Errorf("Expected prev hash to be all zeros, got %x", prevHashStr)
		}
	})

	t.Run("Timestamp", func(t *testing.T) {
		timestamp := header.Timestamp()
		expectedTimestamp := uint32(1231006505)
		if timestamp != expectedTimestamp {
			t.Errorf("Expected timestamp %d, got %d", expectedTimestamp, timestamp)
		}
	})

	t.Run("Bits", func(t *testing.T) {
		bits := header.Bits()
		expectedBits := uint32(0x1d00ffff)
		if bits != expectedBits {
			t.Errorf("Expected bits 0x%x, got 0x%x", expectedBits, bits)
		}
	})

	t.Run("Version", func(t *testing.T) {
		version := header.Version()
		expectedVersion := int32(1)
		if version != expectedVersion {
			t.Errorf("Expected version %d, got %d", expectedVersion, version)
		}
	})

	t.Run("Nonce", func(t *testing.T) {
		nonce := header.Nonce()
		expectedNonce := uint32(2083236893)
		if nonce != expectedNonce {
			t.Errorf("Expected nonce %d, got %d", expectedNonce, nonce)
		}
	})
}

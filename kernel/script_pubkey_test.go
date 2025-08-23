package kernel

import (
	"encoding/hex"
	"errors"
	"testing"
)

func TestInvalidScriptPubkeyData(t *testing.T) {
	// Test with empty data
	_, err := NewScriptPubkeyFromRaw([]byte{})
	if !errors.Is(err, ErrEmptyScriptPubkeyData) {
		t.Errorf("Expected ErrEmptyScriptPubkeyData, got %v", err)
	}
}

func TestScriptPubkeyFromRaw(t *testing.T) {
	scriptHex := "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe26158088ac"
	scriptBytes, err := hex.DecodeString(scriptHex)
	if err != nil {
		t.Fatalf("Failed to decode script hex: %v", err)
	}

	scriptPubkey, err := NewScriptPubkeyFromRaw(scriptBytes)
	if err != nil {
		t.Fatalf("NewScriptPubkeyFromRaw() error = %v", err)
	}
	defer scriptPubkey.Destroy()

	// Test getting the serialized script pubkey
	data, err := scriptPubkey.Bytes()
	if err != nil {
		t.Fatalf("ScriptPubkey.Bytes() error = %v", err)
	}

	if len(data) != len(scriptBytes) {
		t.Errorf("Expected data length %d, got %d", len(scriptBytes), len(data))
	}

	hexStr := hex.EncodeToString(data)
	if hexStr != scriptHex {
		t.Errorf("Expected data hex: %s, got %s", scriptHex, hexStr)
	}
}

func TestScriptPubkeyCopy(t *testing.T) {
	scriptHex := "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe26158088ac"
	scriptBytes, err := hex.DecodeString(scriptHex)
	if err != nil {
		t.Fatalf("Failed to decode script hex: %v", err)
	}

	scriptPubkey, err := NewScriptPubkeyFromRaw(scriptBytes)
	if err != nil {
		t.Fatalf("NewScriptPubkeyFromRaw() error = %v", err)
	}
	defer scriptPubkey.Destroy()

	// Test copying script pubkey
	scriptCopy, err := scriptPubkey.Copy()
	if err != nil {
		t.Fatalf("ScriptPubkey.Copy() error = %v", err)
	}
	if scriptCopy == nil {
		t.Fatal("Copied script pubkey is nil")
	}
	defer scriptCopy.Destroy()

	if scriptCopy.ptr == nil {
		t.Error("Copied script pubkey pointer is nil")
	}

	// Verify copy has same data
	originalData, err := scriptPubkey.Bytes()
	if err != nil {
		t.Fatalf("Original ScriptPubkey.Bytes() error = %v", err)
	}

	copyData, err := scriptCopy.Bytes()
	if err != nil {
		t.Fatalf("Copied ScriptPubkey.Bytes() error = %v", err)
	}

	if hex.EncodeToString(originalData) != hex.EncodeToString(copyData) {
		t.Error("Copied script pubkey data doesn't match original")
	}
}

func TestScriptPubkeyBytes(t *testing.T) {
	scriptHex := "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe26158088ac"
	scriptBytes, err := hex.DecodeString(scriptHex)
	if err != nil {
		t.Fatalf("Failed to decode script hex: %v", err)
	}

	scriptPubkey, err := NewScriptPubkeyFromRaw(scriptBytes)
	if err != nil {
		t.Fatalf("NewScriptPubkeyFromRaw() error = %v", err)
	}
	defer scriptPubkey.Destroy()

	// Test serializing script to bytes
	serialized, err := scriptPubkey.Bytes()
	if err != nil {
		t.Fatalf("ScriptPubkey.Bytes() error = %v", err)
	}

	if len(serialized) == 0 {
		t.Error("Serialized script is empty")
	}

	// The serialized bytes should match the original
	if hex.EncodeToString(serialized) != scriptHex {
		t.Errorf("Serialized script doesn't match original.\nExpected: %s\nGot: %s", scriptHex, hex.EncodeToString(serialized))
	}
}

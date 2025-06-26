package kernel

import (
	"encoding/hex"
	"errors"
	"testing"
)

func TestInvalidScriptPubkeyData(t *testing.T) {
	// Test with empty data
	_, err := NewScriptPubkeyFromRaw([]byte{})
	if !errors.Is(err, ErrInvalidScriptPubkeyData) {
		t.Errorf("Expected ErrInvalidScriptPubkeyData, got %v", err)
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

	// Test getting script pubkey data
	data, err := scriptPubkey.Data()
	if err != nil {
		t.Fatalf("ScriptPubkey.Data() error = %v", err)
	}

	if len(data) != len(scriptBytes) {
		t.Errorf("Expected data length %d, got %d", len(scriptBytes), len(data))
	}

	hexStr := hex.EncodeToString(data)
	if hexStr != scriptHex {
		t.Errorf("Expected data hex: %s, got %s", scriptHex, hexStr)
	}
}

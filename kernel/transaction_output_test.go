package kernel

import (
	"encoding/hex"
	"testing"
)

func TestTransactionOutputCreation(t *testing.T) {
	scriptHex := "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe26158088ac"
	scriptBytes, err := hex.DecodeString(scriptHex)
	if err != nil {
		t.Fatalf("Failed to decode script hex: %v", err)
	}

	scriptPubkey, err := NewScriptPubkey(scriptBytes)
	if err != nil {
		t.Fatalf("Failed to create script pubkey: %v", err)
	}
	defer scriptPubkey.Destroy()

	amount := int64(5000000000)
	output := NewTransactionOutput(scriptPubkey, amount)
	defer output.Destroy()

	gotAmount := output.Amount()
	if gotAmount != amount {
		t.Errorf("Expected amount %d, got %d", amount, gotAmount)
	}

	// Test getting script pubkey
	gotScript := output.ScriptPubkey()

	scriptData, err := gotScript.Bytes()
	if err != nil {
		t.Fatalf("ScriptPubkey.Bytes() error = %v", err)
	}

	if len(scriptData) != len(scriptBytes) {
		t.Errorf("Expected script length %d, got %d", len(scriptBytes), len(scriptData))
	}

	scriptHexGot := hex.EncodeToString(scriptData)
	if scriptHexGot != scriptHex {
		t.Errorf("Expected script hex: %s, got %s", scriptHex, scriptHexGot)
	}
}

func TestTransactionOutputCopy(t *testing.T) {
	scriptHex := "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe26158088ac"
	scriptBytes, err := hex.DecodeString(scriptHex)
	if err != nil {
		t.Fatalf("Failed to decode script hex: %v", err)
	}

	scriptPubkey, err := NewScriptPubkey(scriptBytes)
	if err != nil {
		t.Fatalf("Failed to create script pubkey: %v", err)
	}
	defer scriptPubkey.Destroy()

	amount := int64(5000000000)
	output := NewTransactionOutput(scriptPubkey, amount)
	defer output.Destroy()

	// Test copying transaction output
	outputCopy := output.Copy()
	if outputCopy == nil {
		t.Fatal("Copied transaction output is nil")
	}
	defer outputCopy.Destroy()

	if outputCopy.handle.ptr == nil {
		t.Error("Copied transaction output pointer is nil")
	}

	// Verify copy has same amount
	originalAmount := output.Amount()
	copyAmount := outputCopy.Amount()
	if originalAmount != copyAmount {
		t.Errorf("Copied amount doesn't match: original %d, copy %d", originalAmount, copyAmount)
	}
}

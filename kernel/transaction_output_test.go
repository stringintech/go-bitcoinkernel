package kernel

import (
	"encoding/hex"
	"errors"
	"testing"
)

func TestInvalidTransactionOutput(t *testing.T) {
	_, err := NewTransactionOutput(nil, 1000)
	if !errors.Is(err, ErrScriptPubkeyUninitialized) {
		t.Errorf("Expected ErrScriptPubkeyUninitialized, got %v", err)
	}

	_, err = NewTransactionOutput(&ScriptPubkey{ptr: nil}, 1000)
	if !errors.Is(err, ErrScriptPubkeyUninitialized) {
		t.Errorf("Expected ErrScriptPubkeyUninitialized, got %v", err)
	}
}

func TestTransactionOutputCreation(t *testing.T) {
	scriptHex := "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe26158088ac"
	scriptBytes, err := hex.DecodeString(scriptHex)
	if err != nil {
		t.Fatalf("Failed to decode script hex: %v", err)
	}

	scriptPubkey, err := NewScriptPubkeyFromRaw(scriptBytes)
	if err != nil {
		t.Fatalf("Failed to create script pubkey: %v", err)
	}
	defer scriptPubkey.Destroy()

	amount := int64(5000000000)
	output, err := NewTransactionOutput(scriptPubkey, amount)
	if err != nil {
		t.Fatalf("NewTransactionOutput() error = %v", err)
	}
	defer output.Destroy()

	gotAmount := output.Amount()
	if gotAmount != amount {
		t.Errorf("Expected amount %d, got %d", amount, gotAmount)
	}

	// Test getting script pubkey
	gotScript, err := output.ScriptPubkey()
	if err != nil {
		t.Fatalf("TransactionOutput.ScriptPubkey() error = %v", err)
	}
	defer gotScript.Destroy()

	scriptData, err := gotScript.Data()
	if err != nil {
		t.Fatalf("ScriptPubkey.Data() error = %v", err)
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

	scriptPubkey, err := NewScriptPubkeyFromRaw(scriptBytes)
	if err != nil {
		t.Fatalf("Failed to create script pubkey: %v", err)
	}
	defer scriptPubkey.Destroy()

	amount := int64(5000000000)
	output, err := NewTransactionOutput(scriptPubkey, amount)
	if err != nil {
		t.Fatalf("NewTransactionOutput() error = %v", err)
	}
	defer output.Destroy()

	// Test copying transaction output
	outputCopy, err := output.Copy()
	if err != nil {
		t.Fatalf("TransactionOutput.Copy() error = %v", err)
	}
	if outputCopy == nil {
		t.Fatal("Copied transaction output is nil")
	}
	defer outputCopy.Destroy()

	if outputCopy.ptr == nil {
		t.Error("Copied transaction output pointer is nil")
	}

	// Verify copy has same amount
	originalAmount := output.Amount()
	copyAmount := outputCopy.Amount()
	if originalAmount != copyAmount {
		t.Errorf("Copied amount doesn't match: original %d, copy %d", originalAmount, copyAmount)
	}
}

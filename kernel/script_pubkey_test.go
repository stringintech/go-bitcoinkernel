package kernel

import (
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
)

func TestScriptPubkeyFromRaw(t *testing.T) {
	tests := []struct {
		name      string
		scriptHex string
	}{
		{
			name:      "standard_p2pkh",
			scriptHex: "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe26158088ac",
		},
		{
			name:      "empty_script",
			scriptHex: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scriptBytes, err := hex.DecodeString(tt.scriptHex)
			if err != nil {
				t.Fatalf("Failed to decode script hex: %v", err)
			}

			scriptPubkey := NewScriptPubkey(scriptBytes)
			defer scriptPubkey.Destroy()

			data, err := scriptPubkey.Bytes()
			if err != nil {
				t.Fatalf("ScriptPubkey.Bytes() error = %v", err)
			}

			if len(data) != len(scriptBytes) {
				t.Errorf("Expected data length %d, got %d", len(scriptBytes), len(data))
			}

			hexStr := hex.EncodeToString(data)
			if hexStr != tt.scriptHex {
				t.Errorf("Expected data hex: %s, got %s", tt.scriptHex, hexStr)
			}
		})
	}
}

func TestScriptPubkeyCopy(t *testing.T) {
	scriptHex := "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe26158088ac"
	scriptBytes, err := hex.DecodeString(scriptHex)
	if err != nil {
		t.Fatalf("Failed to decode script hex: %v", err)
	}

	scriptPubkey := NewScriptPubkey(scriptBytes)
	defer scriptPubkey.Destroy()

	// Test copying script pubkey
	scriptCopy := scriptPubkey.Copy()
	if scriptCopy == nil {
		t.Fatal("Copied script pubkey is nil")
	}
	defer scriptCopy.Destroy()

	if scriptCopy.handle.ptr == nil {
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

	scriptPubkey := NewScriptPubkey(scriptBytes)
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

func TestValidScripts(t *testing.T) {
	tests := []struct {
		name            string
		scriptPubkeyHex string
		amount          int64
		txToHex         string
		inputIndex      uint
		spentOutputs    []struct {
			scriptHex string
			amount    int64
		}
		flags       ScriptFlags
		description string
	}{
		{
			name:            "old_style_transaction",
			scriptPubkeyHex: "76a9144bfbaf6afb76cc5771bc6404810d1cc041a6933988ac",
			amount:          0,
			txToHex:         "02000000013f7cebd65c27431a90bba7f796914fe8cc2ddfc3f2cbd6f7e5f2fc854534da95000000006b483045022100de1ac3bcdfb0332207c4a91f3832bd2c2915840165f876ab47c5f8996b971c3602201c6c053d750fadde599e6f5c4e1963df0f01fc0d97815e8157e3d59fe09ca30d012103699b464d1d8bc9e47d4fb1cdaa89a1c5783d68363c4dbc4b524ed3d857148617feffffff02836d3c01000000001976a914fc25d6d5c94003bf5b0c7b640a248e2c637fcfb088ac7ada8202000000001976a914fbed3d9b11183209a57999d54d59f67c019e756c88ac6acb0700",
			inputIndex:      0,
			flags:           ScriptFlagsVerifyAll &^ ScriptFlagsVerifyTaproot,
			description:     "a random old-style transaction from the blockchain",
		},
		{
			name:            "segwit_p2sh_transaction",
			scriptPubkeyHex: "a91434c06f8c87e355e123bdc6dda4ffabc64b6989ef87",
			amount:          1900000,
			txToHex:         "01000000000101d9fd94d0ff0026d307c994d0003180a5f248146efb6371d040c5973f5f66d9df0400000017160014b31b31a6cb654cfab3c50567bcf124f48a0beaecffffffff012cbd1c000000000017a914233b74bf0823fa58bbbd26dfc3bb4ae715547167870247304402206f60569cac136c114a58aedd80f6fa1c51b49093e7af883e605c212bdafcd8d202200e91a55f408a021ad2631bc29a67bd6915b2d7e9ef0265627eabd7f7234455f6012103e7e802f50344303c76d12c089c8724c1b230e3b745693bbe16aad536293d15e300000000",
			inputIndex:      0,
			flags:           ScriptFlagsVerifyAll &^ ScriptFlagsVerifyTaproot,
			description:     "a random segwit transaction from the blockchain using P2SH",
		},
		{
			name:            "native_segwit_transaction",
			scriptPubkeyHex: "0020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d",
			amount:          18393430,
			txToHex:         "010000000001011f97548fbbe7a0db7588a66e18d803d0089315aa7d4cc28360b6ec50ef36718a0100000000ffffffff02df1776000000000017a9146c002a686959067f4866b8fb493ad7970290ab728757d29f0000000000220020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d04004730440220565d170eed95ff95027a69b313758450ba84a01224e1f7f130dda46e94d13f8602207bdd20e307f062594022f12ed5017bbf4a055a06aea91c10110a0e3bb23117fc014730440220647d2dc5b15f60bc37dc42618a370b2a1490293f9e5c8464f53ec4fe1dfe067302203598773895b4b16d37485cbe21b337f4e4b650739880098c592553add7dd4355016952210375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c2103a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff2103c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f88053ae00000000",
			inputIndex:      0,
			flags:           ScriptFlagsVerifyAll &^ ScriptFlagsVerifyTaproot,
			description:     "a random segwit transaction from the blockchain using native segwit",
		},
		{
			name:            "taproot_single_input",
			scriptPubkeyHex: "5120339ce7e165e67d93adb3fef88a6d4beed33f01fa876f05a225242b82a631abc0",
			amount:          88480,
			txToHex:         "01000000000101d1f1c1f8cdf6759167b90f52c9ad358a369f95284e841d7a2536cef31c0549580100000000fdffffff020000000000000000316a2f49206c696b65205363686e6f7272207369677320616e6420492063616e6e6f74206c69652e204062697462756734329e06010000000000225120a37c3903c8d0db6512e2b40b0dffa05e5a3ab73603ce8c9c4b7771e5412328f90140a60c383f71bac0ec919b1d7dbc3eb72dd56e7aa99583615564f9f99b8ae4e837b758773a5b2e4c51348854c8389f008e05029db7f464a5ff2e01d5e6e626174affd30a00",
			inputIndex:      0,
			spentOutputs: []struct {
				scriptHex string
				amount    int64
			}{
				{
					scriptHex: "5120339ce7e165e67d93adb3fef88a6d4beed33f01fa876f05a225242b82a631abc0",
					amount:    88480,
				},
			},
			flags:       ScriptFlagsVerifyAll,
			description: "Single-input Taproot transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := testVerifyScript(t, tt.scriptPubkeyHex, tt.amount, tt.txToHex, tt.inputIndex, tt.flags, tt.spentOutputs)
			if err != nil {
				t.Errorf("testVerifyScript() error = %v", err)
			} else if !valid {
				t.Errorf("testVerifyScript() expected valid script, got invalid")
			}
		})
	}
}

func TestValidTaprootMultiInput(t *testing.T) {
	txToHex := "02000000000102c0f01ead18750892c84b1d4f595149ad38f16847df1fbf490e235b3b78c1f98a0100000000ffffffff456764a19c2682bf5b1567119f06a421849ad1664cf42b5ef95b69d6e2159e9d0000000000ffffffff022202000000000000225120b6c0c2a8ee25a2ae0322ab7f1a06f01746f81f6b90d179c3c2a51a356e6188f1d70e020000000000225120b7da80f57e36930b0515eb09293e25858d13e6b91fee6184943f5a584cb4248e0141933fdc49eb1af1f08ed1e9cf5559259309a8acd25ff1e6999b6955124438aef4fceaa4e6a5f85286631e24837329563595bc3cf4b31e1c687442abb01c4206818101401c9620faf1e8c84187762ad14d04ae3857f59a2f03f1dcbb99290e16dfc572a63b4ea435780a5787af59beb5742fd71cda8a95381517a1ff14b4c67996c4bf8100000000"

	txToBytes, err := hex.DecodeString(txToHex)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	txTo, err := NewTransaction(txToBytes)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}
	defer txTo.Destroy()

	spentOutputs := []struct {
		scriptHex string
		amount    int64
	}{
		{scriptHex: "5120b7da80f57e36930b0515eb09293e25858d13e6b91fee6184943f5a584cb4248e", amount: 546},
		{scriptHex: "5120ab78e077d062e7b8acd7063668b4db5355a1b5d5fd2a46a8e98e62e5e63fab77", amount: 135125},
	}

	outputs := make([]TransactionOutputLike, len(spentOutputs))
	for i, spent := range spentOutputs {
		spentScriptBytes, err := hex.DecodeString(spent.scriptHex)
		if err != nil {
			t.Fatalf("Failed to decode spent script hex: %v", err)
		}
		spentScript := NewScriptPubkey(spentScriptBytes)
		defer spentScript.Destroy()

		output := NewTransactionOutput(spentScript, spent.amount)
		defer output.Destroy()
		outputs[i] = output
	}

	precomputed, err := NewPrecomputedTransactionData(txTo, outputs)
	if err != nil {
		t.Fatalf("Failed to create precomputed transaction data: %v", err)
	}
	defer precomputed.Destroy()

	inputs := []struct {
		scriptPubkeyHex string
		amount          int64
	}{
		{
			scriptPubkeyHex: "5120b7da80f57e36930b0515eb09293e25858d13e6b91fee6184943f5a584cb4248e",
			amount:          546,
		},
		{
			scriptPubkeyHex: "5120ab78e077d062e7b8acd7063668b4db5355a1b5d5fd2a46a8e98e62e5e63fab77",
			amount:          135125,
		},
	}

	for i, input := range inputs {
		t.Run(fmt.Sprintf("input_%d", i), func(t *testing.T) {
			scriptPubkeyBytes, err := hex.DecodeString(input.scriptPubkeyHex)
			if err != nil {
				t.Fatalf("Failed to decode script pubkey hex: %v", err)
			}

			scriptPubkey := NewScriptPubkey(scriptPubkeyBytes)
			defer scriptPubkey.Destroy()

			valid, err := scriptPubkey.Verify(input.amount, txTo, precomputed, uint(i), ScriptFlagsVerifyAll)
			if err != nil {
				t.Errorf("Verify() error = %v", err)
			} else if !valid {
				t.Errorf("Verify() expected valid script, got invalid")
			}
		})
	}
}

func TestInvalidScripts(t *testing.T) {
	tests := []struct {
		name            string
		scriptPubkeyHex string
		amount          int64
		txToHex         string
		inputIndex      uint
		spentOutputs    []struct {
			scriptHex string
			amount    int64
		}
		flags       ScriptFlags
		description string
	}{
		{
			name:            "old_style_wrong_signature",
			scriptPubkeyHex: "76a9144bfbaf6afb76cc5771bc6404810d1cc041a6933988ff", // Modified last byte from 'ac' to 'ff'
			amount:          0,
			txToHex:         "02000000013f7cebd65c27431a90bba7f796914fe8cc2ddfc3f2cbd6f7e5f2fc854534da95000000006b483045022100de1ac3bcdfb0332207c4a91f3832bd2c2915840165f876ab47c5f8996b971c3602201c6c053d750fadde599e6f5c4e1963df0f01fc0d97815e8157e3d59fe09ca30d012103699b464d1d8bc9e47d4fb1cdaa89a1c5783d68363c4dbc4b524ed3d857148617feffffff02836d3c01000000001976a914fc25d6d5c94003bf5b0c7b640a248e2c637fcfb088ac7ada8202000000001976a914fbed3d9b11183209a57999d54d59f67c019e756c88ac6acb0700",
			inputIndex:      0,
			flags:           ScriptFlagsVerifyAll & ^ScriptFlagsVerifyTaproot,
			description:     "a random old-style transaction from the blockchain - WITH WRONG SIGNATURE for the address",
		},
		{
			name:            "segwit_p2sh_wrong_amount",
			scriptPubkeyHex: "a91434c06f8c87e355e123bdc6dda4ffabc64b6989ef87",
			amount:          900000, // Wrong amount, should be 1900000
			txToHex:         "01000000000101d9fd94d0ff0026d307c994d0003180a5f248146efb6371d040c5973f5f66d9df0400000017160014b31b31a6cb654cfab3c50567bcf124f48a0beaecffffffff012cbd1c000000000017a914233b74bf0823fa58bbbd26dfc3bb4ae715547167870247304402206f60569cac136c114a58aedd80f6fa1c51b49093e7af883e605c212bdafcd8d202200e91a55f408a021ad2631bc29a67bd6915b2d7e9ef0265627eabd7f7234455f6012103e7e802f50344303c76d12c089c8724c1b230e3b745693bbe16aad536293d15e300000000",
			inputIndex:      0,
			flags:           ScriptFlagsVerifyAll & ^ScriptFlagsVerifyTaproot,
			description:     "a random segwit transaction from the blockchain using P2SH - WITH WRONG AMOUNT",
		},
		{
			name:            "native_segwit_wrong_segwit",
			scriptPubkeyHex: "0020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58f", // Modified last byte from 'd' to 'f'
			amount:          18393430,
			txToHex:         "010000000001011f97548fbbe7a0db7588a66e18d803d0089315aa7d4cc28360b6ec50ef36718a0100000000ffffffff02df1776000000000017a9146c002a686959067f4866b8fb493ad7970290ab728757d29f0000000000220020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d04004730440220565d170eed95ff95027a69b313758450ba84a01224e1f7f130dda46e94d13f8602207bdd20e307f062594022f12ed5017bbf4a055a06aea91c10110a0e3bb23117fc014730440220647d2dc5b15f60bc37dc42618a370b2a1490293f9e5c8464f53ec4fe1dfe067302203598773895b4b16d37485cbe21b337f4e4b650739880098c592553add7dd4355016952210375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c2103a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff2103c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f88053ae00000000",
			inputIndex:      0,
			flags:           ScriptFlagsVerifyAll & ^ScriptFlagsVerifyTaproot,
			description:     "a random segwit transaction from the blockchain using native segwit - WITH WRONG SEGWIT",
		},
		{
			name:            "empty_scriptpubkey",
			scriptPubkeyHex: "",
			amount:          0,
			// Minimal coinbase-style transaction with a single empty scriptSig and zero-value output;
			// used to trigger verification paths.
			txToHex:     "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff00ffffffff0100000000000000000000000000",
			inputIndex:  0,
			flags:       ScriptFlagsVerifyAll & ^ScriptFlagsVerifyTaproot,
			description: "empty scriptPubkey should fail verification",
		},
		{
			name:            "taproot_wrong_spent_output_amount",
			scriptPubkeyHex: "5120339ce7e165e67d93adb3fef88a6d4beed33f01fa876f05a225242b82a631abc0",
			amount:          88480,
			txToHex:         "01000000000101d1f1c1f8cdf6759167b90f52c9ad358a369f95284e841d7a2536cef31c0549580100000000fdffffff020000000000000000316a2f49206c696b65205363686e6f7272207369677320616e6420492063616e6e6f74206c69652e204062697462756734329e06010000000000225120a37c3903c8d0db6512e2b40b0dffa05e5a3ab73603ce8c9c4b7771e5412328f90140a60c383f71bac0ec919b1d7dbc3eb72dd56e7aa99583615564f9f99b8ae4e837b758773a5b2e4c51348854c8389f008e05029db7f464a5ff2e01d5e6e626174affd30a00",
			inputIndex:      0,
			spentOutputs: []struct {
				scriptHex string
				amount    int64
			}{
				{
					scriptHex: "5120339ce7e165e67d93adb3fef88a6d4beed33f01fa876f05a225242b82a631abc0",
					amount:    100000, // Wrong amount, should be 88480
				},
			},
			flags:       ScriptFlagsVerifyAll,
			description: "Taproot transaction with incorrect spent output amount in precomputed data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := testVerifyScript(t, tt.scriptPubkeyHex, tt.amount, tt.txToHex, tt.inputIndex, tt.flags, tt.spentOutputs)
			if err != nil {
				t.Errorf("testVerifyScript() unexpected error = %v", err)
			} else if valid {
				t.Errorf("testVerifyScript() expected invalid script, got valid")
			}
		})
	}
}

func TestScriptVerifyErrors(t *testing.T) {
	// Use a valid transaction for testing error conditions
	validScriptHex := "76a9144bfbaf6afb76cc5771bc6404810d1cc041a6933988ac"
	validTxHex := "02000000013f7cebd65c27431a90bba7f796914fe8cc2ddfc3f2cbd6f7e5f2fc854534da95000000006b483045022100de1ac3bcdfb0332207c4a91f3832bd2c2915840165f876ab47c5f8996b971c3602201c6c053d750fadde599e6f5c4e1963df0f01fc0d97815e8157e3d59fe09ca30d012103699b464d1d8bc9e47d4fb1cdaa89a1c5783d68363c4dbc4b524ed3d857148617feffffff02836d3c01000000001976a914fc25d6d5c94003bf5b0c7b640a248e2c637fcfb088ac7ada8202000000001976a914fbed3d9b11183209a57999d54d59f67c019e756c88ac6acb0700"

	scriptBytes, err := hex.DecodeString(validScriptHex)
	if err != nil {
		t.Fatalf("Failed to decode script hex: %v", err)
	}

	scriptPubkey := NewScriptPubkey(scriptBytes)
	defer scriptPubkey.Destroy()

	txBytes, err := hex.DecodeString(validTxHex)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	tx, err := NewTransaction(txBytes)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}
	defer tx.Destroy()

	// Create a precomputed txdata without spent outputs for the spent_outputs_required test
	emptyPrecomputedTxData, err := NewPrecomputedTransactionData(tx, nil)
	if err != nil {
		t.Fatalf("Failed to create empty precomputed transaction data: %v", err)
	}
	defer emptyPrecomputedTxData.Destroy()

	tests := []struct {
		name              string
		inputIndex        uint
		flags             ScriptFlags
		precomputedTxData *PrecomputedTransactionData
		expectedError     error
		description       string
	}{
		{
			name:              "invalid_input_index",
			inputIndex:        999,
			flags:             ScriptFlagsVerifyAll,
			precomputedTxData: nil,
			expectedError:     ErrVerifyScriptVerifyTxInputIndex,
			description:       "input index out of bounds should return error",
		},
		{
			name:              "invalid_flags",
			inputIndex:        0,
			flags:             0xFFFFFFFF, // Invalid flags
			precomputedTxData: nil,
			expectedError:     ErrVerifyScriptVerifyInvalidFlags,
			description:       "invalid flags should return error",
		},
		{
			name:              "spent_outputs_required",
			inputIndex:        0,
			flags:             ScriptFlagsVerifyAll,
			precomputedTxData: emptyPrecomputedTxData,
			expectedError:     ErrVerifyScriptVerifySpentOutputsRequired,
			description:       "taproot verification with empty spent outputs should return error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := scriptPubkey.Verify(0, tx, tt.precomputedTxData, tt.inputIndex, tt.flags)
			if err == nil {
				t.Errorf("Expected error %v, got nil (valid=%v)", tt.expectedError, valid)
				return
			}
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

// testVerifyScript is a helper function that creates necessary objects and calls VerifyScript
func testVerifyScript(
	t *testing.T,
	scriptPubkeyHex string,
	amount int64,
	txToHex string,
	inputIndex uint,
	flags ScriptFlags,
	spentOutputs []struct {
		scriptHex string
		amount    int64
	},
) (bool, error) {
	scriptPubkeyBytes, err := hex.DecodeString(scriptPubkeyHex)
	if err != nil {
		t.Fatalf("Failed to decode script pubkey hex: %v", err)
	}

	scriptPubkey := NewScriptPubkey(scriptPubkeyBytes)
	defer scriptPubkey.Destroy()

	txToBytes, err := hex.DecodeString(txToHex)
	if err != nil {
		t.Fatalf("Failed to decode transaction hex: %v", err)
	}

	txTo, err := NewTransaction(txToBytes)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}
	defer txTo.Destroy()

	// Create spent outputs if provided
	var outputs []TransactionOutputLike
	if len(spentOutputs) > 0 {
		outputs = make([]TransactionOutputLike, len(spentOutputs))
		for i, spent := range spentOutputs {
			spentScriptBytes, err := hex.DecodeString(spent.scriptHex)
			if err != nil {
				t.Fatalf("Failed to decode spent script hex: %v", err)
			}
			spentScript := NewScriptPubkey(spentScriptBytes)
			defer spentScript.Destroy()

			output := NewTransactionOutput(spentScript, spent.amount)
			defer output.Destroy()
			outputs[i] = output
		}
	}

	// Always create precomputed data (outputs can be nil)
	precomputed, err := NewPrecomputedTransactionData(txTo, outputs)
	if err != nil {
		t.Fatalf("Failed to create precomputed transaction data: %v", err)
	}
	defer precomputed.Destroy()

	return scriptPubkey.Verify(amount, txTo, precomputed, inputIndex, flags)
}

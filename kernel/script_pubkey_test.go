package kernel

import (
	"encoding/hex"
	"errors"
	"testing"
)

func TestScriptPubkeyFromRaw(t *testing.T) {
	scriptHex := "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe26158088ac"
	scriptBytes, err := hex.DecodeString(scriptHex)
	if err != nil {
		t.Fatalf("Failed to decode script hex: %v", err)
	}

	scriptPubkey, err := NewScriptPubkey(scriptBytes)
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

	scriptPubkey, err := NewScriptPubkey(scriptBytes)
	if err != nil {
		t.Fatalf("NewScriptPubkeyFromRaw() error = %v", err)
	}
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

	scriptPubkey, err := NewScriptPubkey(scriptBytes)
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

func TestValidScripts(t *testing.T) {
	tests := []struct {
		name            string
		scriptPubkeyHex string
		amount          int64
		txToHex         string
		inputIndex      uint
		description     string
	}{
		{
			name:            "old_style_transaction",
			scriptPubkeyHex: "76a9144bfbaf6afb76cc5771bc6404810d1cc041a6933988ac",
			amount:          0,
			txToHex:         "02000000013f7cebd65c27431a90bba7f796914fe8cc2ddfc3f2cbd6f7e5f2fc854534da95000000006b483045022100de1ac3bcdfb0332207c4a91f3832bd2c2915840165f876ab47c5f8996b971c3602201c6c053d750fadde599e6f5c4e1963df0f01fc0d97815e8157e3d59fe09ca30d012103699b464d1d8bc9e47d4fb1cdaa89a1c5783d68363c4dbc4b524ed3d857148617feffffff02836d3c01000000001976a914fc25d6d5c94003bf5b0c7b640a248e2c637fcfb088ac7ada8202000000001976a914fbed3d9b11183209a57999d54d59f67c019e756c88ac6acb0700",
			inputIndex:      0,
			description:     "a random old-style transaction from the blockchain",
		},
		{
			name:            "segwit_p2sh_transaction",
			scriptPubkeyHex: "a91434c06f8c87e355e123bdc6dda4ffabc64b6989ef87",
			amount:          1900000,
			txToHex:         "01000000000101d9fd94d0ff0026d307c994d0003180a5f248146efb6371d040c5973f5f66d9df0400000017160014b31b31a6cb654cfab3c50567bcf124f48a0beaecffffffff012cbd1c000000000017a914233b74bf0823fa58bbbd26dfc3bb4ae715547167870247304402206f60569cac136c114a58aedd80f6fa1c51b49093e7af883e605c212bdafcd8d202200e91a55f408a021ad2631bc29a67bd6915b2d7e9ef0265627eabd7f7234455f6012103e7e802f50344303c76d12c089c8724c1b230e3b745693bbe16aad536293d15e300000000",
			inputIndex:      0,
			description:     "a random segwit transaction from the blockchain using P2SH",
		},
		{
			name:            "native_segwit_transaction",
			scriptPubkeyHex: "0020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d",
			amount:          18393430,
			txToHex:         "010000000001011f97548fbbe7a0db7588a66e18d803d0089315aa7d4cc28360b6ec50ef36718a0100000000ffffffff02df1776000000000017a9146c002a686959067f4866b8fb493ad7970290ab728757d29f0000000000220020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d04004730440220565d170eed95ff95027a69b313758450ba84a01224e1f7f130dda46e94d13f8602207bdd20e307f062594022f12ed5017bbf4a055a06aea91c10110a0e3bb23117fc014730440220647d2dc5b15f60bc37dc42618a370b2a1490293f9e5c8464f53ec4fe1dfe067302203598773895b4b16d37485cbe21b337f4e4b650739880098c592553add7dd4355016952210375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c2103a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff2103c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f88053ae00000000",
			inputIndex:      0,
			description:     "a random segwit transaction from the blockchain using native segwit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testVerifyScript(t, tt.scriptPubkeyHex, tt.amount, tt.txToHex, tt.inputIndex)
			if err != nil {
				t.Errorf("testVerifyScript() error = %v", err)
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
		description     string
	}{
		{
			name:            "old_style_wrong_signature",
			scriptPubkeyHex: "76a9144bfbaf6afb76cc5771bc6404810d1cc041a6933988ff", // Modified last byte from 'ac' to 'ff'
			amount:          0,
			txToHex:         "02000000013f7cebd65c27431a90bba7f796914fe8cc2ddfc3f2cbd6f7e5f2fc854534da95000000006b483045022100de1ac3bcdfb0332207c4a91f3832bd2c2915840165f876ab47c5f8996b971c3602201c6c053d750fadde599e6f5c4e1963df0f01fc0d97815e8157e3d59fe09ca30d012103699b464d1d8bc9e47d4fb1cdaa89a1c5783d68363c4dbc4b524ed3d857148617feffffff02836d3c01000000001976a914fc25d6d5c94003bf5b0c7b640a248e2c637fcfb088ac7ada8202000000001976a914fbed3d9b11183209a57999d54d59f67c019e756c88ac6acb0700",
			inputIndex:      0,
			description:     "a random old-style transaction from the blockchain - WITH WRONG SIGNATURE for the address",
		},
		{
			name:            "segwit_p2sh_wrong_amount",
			scriptPubkeyHex: "a91434c06f8c87e355e123bdc6dda4ffabc64b6989ef87",
			amount:          900000, // Wrong amount, should be 1900000
			txToHex:         "01000000000101d9fd94d0ff0026d307c994d0003180a5f248146efb6371d040c5973f5f66d9df0400000017160014b31b31a6cb654cfab3c50567bcf124f48a0beaecffffffff012cbd1c000000000017a914233b74bf0823fa58bbbd26dfc3bb4ae715547167870247304402206f60569cac136c114a58aedd80f6fa1c51b49093e7af883e605c212bdafcd8d202200e91a55f408a021ad2631bc29a67bd6915b2d7e9ef0265627eabd7f7234455f6012103e7e802f50344303c76d12c089c8724c1b230e3b745693bbe16aad536293d15e300000000",
			inputIndex:      0,
			description:     "a random segwit transaction from the blockchain using P2SH - WITH WRONG AMOUNT",
		},
		{
			name:            "native_segwit_wrong_segwit",
			scriptPubkeyHex: "0020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58f", // Modified last byte from 'd' to 'f'
			amount:          18393430,
			txToHex:         "010000000001011f97548fbbe7a0db7588a66e18d803d0089315aa7d4cc28360b6ec50ef36718a0100000000ffffffff02df1776000000000017a9146c002a686959067f4866b8fb493ad7970290ab728757d29f0000000000220020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d04004730440220565d170eed95ff95027a69b313758450ba84a01224e1f7f130dda46e94d13f8602207bdd20e307f062594022f12ed5017bbf4a055a06aea91c10110a0e3bb23117fc014730440220647d2dc5b15f60bc37dc42618a370b2a1490293f9e5c8464f53ec4fe1dfe067302203598773895b4b16d37485cbe21b337f4e4b650739880098c592553add7dd4355016952210375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c2103a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff2103c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f88053ae00000000",
			inputIndex:      0,
			description:     "a random segwit transaction from the blockchain using native segwit - WITH WRONG SEGWIT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var scriptVerifyError *ScriptVerifyError
			err := testVerifyScript(t, tt.scriptPubkeyHex, tt.amount, tt.txToHex, tt.inputIndex)
			if err == nil || !errors.As(err, &scriptVerifyError) {
				t.Errorf("testVerifyScript() was expected to fail with ScriptVerifyError, got: %v", err)
			}
		})
	}
}

// testVerifyScript is a helper function that creates the necessary objects and calls VerifyScript
func testVerifyScript(t *testing.T, scriptPubkeyHex string, amount int64, txToHex string, inputIndex uint) error {
	scriptPubkeyBytes, err := hex.DecodeString(scriptPubkeyHex)
	if err != nil {
		t.Fatalf("Failed to decode script pubkey hex: %v", err)
	}

	scriptPubkey, err := NewScriptPubkey(scriptPubkeyBytes)
	if err != nil {
		t.Fatalf("Failed to create script pubkey: %v", err)
	}
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

	flags := ScriptFlags(ScriptFlagsVerifyAll &^ ScriptFlagsVerifyTaproot)
	return scriptPubkey.Verify(amount, txTo, nil, inputIndex, flags)
}

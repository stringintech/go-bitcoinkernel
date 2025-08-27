package kernel

import (
	"encoding/hex"
	"testing"
)

func TestValidationInterfaceCallbacks(t *testing.T) {
	var blockCheckedCalled bool
	var lastValidationMode ValidationMode
	var lastBlockData []byte

	callbacks := &ValidationInterfaceCallbacks{
		OnBlockChecked: func(block *Block, state *BlockValidationState) {
			blockCheckedCalled = true
			lastValidationMode = state.ValidationMode()
			var err error
			lastBlockData, err = block.Bytes()
			if err != nil {
				t.Fatal(err)
			}
		},
	}

	suite := ChainstateManagerTestSuite{
		MaxBlockHeightToImport: 2,
		ValidationCallbacks:    callbacks,
	}
	suite.Setup(t)

	if !blockCheckedCalled {
		t.Error("OnBlockChecked callback was not called")
	}

	if lastValidationMode != ValidationStateValid {
		t.Errorf("Expected validation mode %d, got %d", ValidationStateValid, lastValidationMode)
	}

	lastBlockDataHex := hex.EncodeToString(lastBlockData)
	if lastBlockDataHex != "00000020a629da61ccd6c9de14dd22d4dcf06ac4b98828801fb58275af1ed2c89e361b79677daedb5fc7781c5907a88133cd461b4865e9a4881fecfb362304ad1806acf3a7242d66ffff7f200100000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025200ffffffff0200f2052a010000001600141409745405c4e8310a875bcd602db6b9b3dc0cf90000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000" {
		t.Errorf("Unexpected block data for last block")
	}
}

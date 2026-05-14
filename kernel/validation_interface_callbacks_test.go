package kernel

import (
	"encoding/hex"
	"testing"
)

type validationTestCallbacks struct {
	lastBlockCheckedData []byte
	lastValidationMode   ValidationMode
	lastConnectedData    []byte
	lastConnectedHeight  int32
}

func (c *validationTestCallbacks) BlockChecked(block *Block, state *BlockValidationStateView) {
	c.lastValidationMode = state.ValidationMode()
	c.lastBlockCheckedData, _ = block.Bytes()
}

func (c *validationTestCallbacks) BlockConnected(block *Block, entry *BlockTreeEntry) {
	c.lastConnectedHeight = entry.Height()
	c.lastConnectedData, _ = block.Bytes()
}

func (c *validationTestCallbacks) PoWValidBlock(_ *Block, _ *BlockTreeEntry)     {}
func (c *validationTestCallbacks) BlockDisconnected(_ *Block, _ *BlockTreeEntry) {}

func TestValidationInterfaceCallbacks(t *testing.T) {
	cb := &validationTestCallbacks{}
	suite := ChainstateManagerTestSuite{
		MaxBlockHeightToImport: 2,
		ValidationCallbacks:    cb,
	}
	suite.Setup(t)

	if cb.lastBlockCheckedData == nil {
		t.Error("BlockChecked callback was not called")
	}
	if cb.lastValidationMode != ValidationStateValid {
		t.Errorf("Expected validation mode %d, got %d", ValidationStateValid, cb.lastValidationMode)
	}
	expectedHex := "00000020a629da61ccd6c9de14dd22d4dcf06ac4b98828801fb58275af1ed2c89e361b79677daedb5fc7781c5907a88133cd461b4865e9a4881fecfb362304ad1806acf3a7242d66ffff7f200100000001020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff025200ffffffff0200f2052a010000001600141409745405c4e8310a875bcd602db6b9b3dc0cf90000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000"
	if hex.EncodeToString(cb.lastBlockCheckedData) != expectedHex {
		t.Errorf("Unexpected block data for last block")
	}
	if cb.lastConnectedData == nil {
		t.Error("BlockConnected callback was not called")
	}
	if cb.lastConnectedHeight != 2 {
		t.Errorf("Expected connected block height 2, got %d", cb.lastConnectedHeight)
	}
	if hex.EncodeToString(cb.lastConnectedData) != expectedHex {
		t.Errorf("Unexpected block data for connected block")
	}
}

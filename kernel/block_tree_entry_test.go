package kernel

import (
	"testing"
)

func TestBlockTreeEntryGetPrevious(t *testing.T) {
	suite := ChainstateManagerTestSuite{
		MaxBlockHeightToImport: 2, // Import just genesis and first block
	}
	suite.Setup(t)

	chain := suite.Manager.GetActiveChain()

	// Get block at height 1
	entry := chain.GetByHeight(1)

	// Test getting previous block (should be genesis)
	prevEntry := entry.Previous()
	if prevEntry == nil {
		t.Fatal("Previous block tree entry is nil")
	}

	// Verify previous block is genesis (height 0)
	previousHeight := prevEntry.Height()
	if previousHeight != 0 {
		t.Errorf("Expected previous block height 0, got %d", previousHeight)
	}

	// Test genesis block has no previous
	genesisEntry := chain.GetByHeight(0)

	// Genesis should have no previous block (should return nil)
	genesisPrevious := genesisEntry.Previous()
	if genesisPrevious != nil {
		t.Error("Genesis block should not have a previous block")
	}
}

func TestBlockTreeEntryEquals(t *testing.T) {
	suite := ChainstateManagerTestSuite{
		MaxBlockHeightToImport: 2,
	}
	suite.Setup(t)

	chain := suite.Manager.GetActiveChain()

	// Same entry should equal itself
	entry1 := chain.GetByHeight(1)
	if !entry1.Equals(entry1) {
		t.Error("Entry should equal itself")
	}

	// Different retrievals of same height should be equal
	entry1Again := chain.GetByHeight(1)
	if !entry1.Equals(entry1Again) {
		t.Error("Same height entries should be equal")
	}

	// Different heights should not be equal
	entry0 := chain.GetByHeight(0)
	if entry1.Equals(entry0) {
		t.Error("Different height entries should not be equal")
	}

	// Nil comparison should return false
	if entry1.Equals(nil) {
		t.Error("Entry should not equal nil")
	}
}

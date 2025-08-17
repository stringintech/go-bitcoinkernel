package kernel

import (
	"testing"
)

func TestBlockTreeEntryGetPrevious(t *testing.T) {
	suite := ChainstateManagerTestSuite{
		MaxBlockHeightToImport: 2, // Import just genesis and first block
	}
	suite.Setup(t)

	chain, err := suite.Manager.GetActiveChain()
	if err != nil {
		t.Fatalf("GetActiveChain() error = %v", err)
	}
	defer chain.Destroy()

	// Get block at height 1
	entry, err := chain.GetByHeight(1)
	if err != nil {
		t.Fatalf("GetByHeight(1) error = %v", err)
	}
	defer entry.Destroy()

	// Test getting previous block (should be genesis)
	prevEntry := entry.Previous()
	if prevEntry == nil {
		t.Fatal("Previous block tree entry is nil")
	}
	defer prevEntry.Destroy()

	// Verify previous block is genesis (height 0)
	previousHeight := prevEntry.Height()
	if previousHeight != 0 {
		t.Errorf("Expected previous block height 0, got %d", previousHeight)
	}

	// Test genesis block has no previous
	genesisEntry, err := chain.GetGenesis()
	if err != nil {
		t.Fatalf("GetGenesis() error = %v", err)
	}
	defer genesisEntry.Destroy()

	// Genesis should have no previous block (should return nil)
	genesisPrevious := genesisEntry.Previous()
	if genesisPrevious != nil {
		t.Error("Genesis block should not have a previous block")
		genesisPrevious.Destroy()
	}
}

package kernel

import (
	"testing"
)

func TestChain(t *testing.T) {
	suite := ChainstateManagerTestSuite{
		MaxBlockHeightToImport: 3,   // Import genesis and first few blocks
		NotificationCallbacks:  nil, // no notification callbacks
		ValidationCallbacks:    nil, // no validation callbacks
	}
	suite.Setup(t)

	chain, err := suite.Manager.GetActiveChain()
	if err != nil {
		t.Fatalf("GetActiveChain() error = %v", err)
	}
	defer chain.Destroy()

	// Test GetGenesis
	genesis, err := chain.GetGenesis()
	if err != nil {
		t.Fatalf("GetGenesis() error = %v", err)
	}
	defer genesis.Destroy()

	height := genesis.Height()
	if height != 0 {
		t.Errorf("Expected genesis height 0, got %d", height)
	}

	genesisHash, err := genesis.Hash()
	if err != nil {
		t.Fatalf("BlockIndex.Hash() error = %v", err)
	}
	defer genesisHash.Destroy()

	hashBytes := genesisHash.Bytes()
	if len(hashBytes) != 32 {
		t.Errorf("Expected hash length 32, got %d", len(hashBytes))
	}

	// Test GetTip
	tip, err := chain.GetTip()
	if err != nil {
		t.Fatalf("GetTip() error = %v", err)
	}
	defer tip.Destroy()

	tipHeight := tip.Height()
	if tipHeight <= 0 {
		t.Errorf("Expected tip height > 0, got %d", tipHeight)
	}

	tipHash, err := tip.Hash()
	if err != nil {
		t.Fatalf("Failed to get tip hash: %v", err)
	}
	defer tipHash.Destroy()

	tipHashBytes := tipHash.Bytes()
	if len(tipHashBytes) != 32 {
		t.Errorf("Expected hash length 32, got %d", len(tipHashBytes))
	}

	if tip.Height() != suite.ImportedBlocksCount {
		t.Errorf("Expected tip height %d, got %d", suite.ImportedBlocksCount, tip.Height())
	}

	// Test GetByHeight
	block1, err := chain.GetByHeight(1)
	if err != nil {
		t.Fatalf("GetByHeight(1) error = %v", err)
	}
	defer block1.Destroy()

	if block1.Height() != 1 {
		t.Errorf("Expected block height 1, got %d", block1.Height())
	}

	// Test GetNextBlockTreeEntry
	nextEntry, err := chain.GetNextBlockTreeEntry(genesis)
	if err != nil {
		t.Fatalf("Chain.GetNextBlockTreeEntry() error = %v", err)
	}
	if nextEntry == nil {
		t.Fatal("Next block tree entry is nil")
	}
	defer nextEntry.Destroy()

	if nextEntry.Height() != 1 {
		t.Errorf("Expected next block height 1, got %d", nextEntry.Height())
	}

	// Test Contains
	containsGenesis := chain.Contains(genesis)
	if !containsGenesis {
		t.Error("Chain should contain genesis block")
	}

	containsBlock1 := chain.Contains(block1)
	if !containsBlock1 {
		t.Error("Chain should contain block at height 1")
	}
}

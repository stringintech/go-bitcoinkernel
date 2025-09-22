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

	chain := suite.Manager.GetActiveChain()

	// Test GetGenesis
	genesis := chain.GetGenesis()
	height := genesis.Height()
	if height != 0 {
		t.Errorf("Expected genesis height 0, got %d", height)
	}

	// Test GetTip
	tip := chain.GetTip()
	tipHeight := tip.Height()
	if tipHeight <= 0 {
		t.Errorf("Expected tip height > 0, got %d", tipHeight)
	}
	if tip.Height() != suite.ImportedBlocksCount {
		t.Errorf("Expected tip height %d, got %d", suite.ImportedBlocksCount, tip.Height())
	}

	// Test GetHeight
	chainHeight := chain.GetHeight()
	if chainHeight != tipHeight {
		t.Errorf("Expected chain height %d to match tip height %d", chainHeight, tipHeight)
	}

	// Test GetByHeight
	block1 := chain.GetByHeight(1)

	if block1.Height() != 1 {
		t.Errorf("Expected block height 1, got %d", block1.Height())
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

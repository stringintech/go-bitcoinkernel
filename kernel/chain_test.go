package kernel

import (
	"slices"
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

	t.Run("GetGenesis", func(t *testing.T) {
		genesis := chain.GetGenesis()
		height := genesis.Height()
		if height != 0 {
			t.Errorf("Expected genesis height 0, got %d", height)
		}
	})

	t.Run("GetTip", func(t *testing.T) {
		tip := chain.GetTip()
		tipHeight := tip.Height()
		if tipHeight <= 0 {
			t.Errorf("Expected tip height > 0, got %d", tipHeight)
		}
		if tip.Height() != suite.ImportedBlocksCount {
			t.Errorf("Expected tip height %d, got %d", suite.ImportedBlocksCount, tip.Height())
		}
	})

	t.Run("GetHeight", func(t *testing.T) {
		chainHeight := chain.GetHeight()
		tipHeight := chain.GetTip().Height()
		if chainHeight != tipHeight {
			t.Errorf("Expected chain height %d to match tip height %d", chainHeight, tipHeight)
		}
	})

	t.Run("GetByHeight", func(t *testing.T) {
		block1 := chain.GetByHeight(1)
		if block1.Height() != 1 {
			t.Errorf("Expected block height 1, got %d", block1.Height())
		}
	})

	t.Run("Contains", func(t *testing.T) {
		genesis := chain.GetGenesis()
		if !chain.Contains(genesis) {
			t.Error("Chain should contain genesis block")
		}

		block1 := chain.GetByHeight(1)
		if !chain.Contains(block1) {
			t.Error("Chain should contain block at height 1")
		}
	})

	t.Run("Entries", func(t *testing.T) {
		count := len(slices.Collect(chain.Entries()))
		expectedCount := int(chain.GetHeight()) + 1 // height 0 to tip inclusive
		if count != expectedCount {
			t.Errorf("Expected to iterate over %d entries, got %d", expectedCount, count)
		}
	})

	t.Run("EntriesRange", func(t *testing.T) {
		count := len(slices.Collect(chain.EntriesRange(0, 2)))
		if count != 2 {
			t.Errorf("Expected to iterate over 2 entries, got %d", count)
		}

		count = len(slices.Collect(chain.EntriesRange(1, 1000)))
		expectedCount := int(chain.GetHeight()) // height 1 to tip inclusive
		if count != expectedCount {
			t.Errorf("Expected to iterate over %d entries, got %d", expectedCount, count)
		}
	})

	t.Run("EntriesFrom", func(t *testing.T) {
		count := len(slices.Collect(chain.EntriesFrom(0)))
		expectedCount := int(chain.GetHeight()) + 1
		if count != expectedCount {
			t.Errorf("Expected to iterate over %d entries, got %d", expectedCount, count)
		}

		count = len(slices.Collect(chain.EntriesFrom(1)))
		expectedCount = int(chain.GetHeight())
		if count != expectedCount {
			t.Errorf("Expected to iterate over %d entries, got %d", expectedCount, count)
		}
	})
}

package kernel

import (
	"slices"
	"testing"
)

func TestChain(t *testing.T) {
	suite := ChainstateManagerTestSuite{
		MaxBlockHeightToImport: 3,
		NotificationCallbacks:  nil,
		ValidationCallbacks:    nil,
	}
	suite.Setup(t)

	chain := suite.Manager.GetActiveChain()

	t.Run("GetHeight", func(t *testing.T) {
		chainHeight := chain.GetHeight()
		if chainHeight != suite.MaxBlockHeightToImport {
			t.Errorf("Expected chain height %d to match tip height %d", chainHeight, suite.MaxBlockHeightToImport)
		}
	})

	t.Run("GetByHeight", func(t *testing.T) {
		block1 := chain.GetByHeight(1)
		if block1.Height() != 1 {
			t.Errorf("Expected block height 1, got %d", block1.Height())
		}
	})

	t.Run("Contains", func(t *testing.T) {
		genesis := chain.GetByHeight(0)
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

package kernel

import (
	"testing"
)

func TestBlockTreeEntry(t *testing.T) {
	suite := ChainstateManagerTestSuite{
		MaxBlockHeightToImport: 2,
	}
	suite.Setup(t)

	chain := suite.Manager.GetActiveChain()

	t.Run("Previous", func(t *testing.T) {
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
	})

	t.Run("Equals", func(t *testing.T) {
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
	})

	t.Run("GetHeader", func(t *testing.T) {
		// Get genesis block entry
		genesisEntry := chain.GetByHeight(0)
		if genesisEntry == nil {
			t.Fatal("Genesis block tree entry is nil")
		}

		// Get the header from the entry
		header := genesisEntry.GetHeader()
		if header == nil {
			t.Fatal("GetHeader() returned nil")
		}
		defer header.Destroy()

		if header.ptr == nil {
			t.Error("Header pointer is nil")
		}

		// Verify the header hash matches the regtest genesis block
		hash := header.Hash()
		defer hash.Destroy()

		expectedHash := "0f9188f13cb7b2c71f2a335e3a4fc328bf5beb436012afca590b1a11466e2206"
		actualHashHex := hash.String()
		if actualHashHex != expectedHash {
			t.Errorf("Expected header hash %s, got %s", expectedHash, actualHashHex)
		}
	})
}

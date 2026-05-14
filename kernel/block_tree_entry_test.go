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
		entry, err := chain.GetByHeight(1)
		if err != nil {
			t.Fatalf("GetByHeight(1) error = %v", err)
		}

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
		genesisEntry, err := chain.GetByHeight(0)
		if err != nil {
			t.Fatalf("GetByHeight(0) error = %v", err)
		}

		// Genesis should have no previous block (should return nil)
		genesisPrevious := genesisEntry.Previous()
		if genesisPrevious != nil {
			t.Error("Genesis block should not have a previous block")
		}
	})

	t.Run("Equals", func(t *testing.T) {
		// Same entry should equal itself
		entry1, err := chain.GetByHeight(1)
		if err != nil {
			t.Fatalf("GetByHeight(1) error = %v", err)
		}
		if !entry1.Equals(entry1) {
			t.Error("Entry should equal itself")
		}

		// Different retrievals of same height should be equal
		entry1Again, err := chain.GetByHeight(1)
		if err != nil {
			t.Fatalf("GetByHeight(1) error = %v", err)
		}
		if !entry1.Equals(entry1Again) {
			t.Error("Same height entries should be equal")
		}

		// Different heights should not be equal
		entry0, err := chain.GetByHeight(0)
		if err != nil {
			t.Fatalf("GetByHeight(0) error = %v", err)
		}
		if entry1.Equals(entry0) {
			t.Error("Different height entries should not be equal")
		}

		// Nil comparison should return false
		if entry1.Equals(nil) {
			t.Error("Entry should not equal nil")
		}
	})

	t.Run("Ancestor", func(t *testing.T) {
		entry2, err := chain.GetByHeight(2)
		if err != nil {
			t.Fatalf("GetByHeight(2) error = %v", err)
		}

		for _, height := range []int32{0, 1, 2} {
			ancestor := entry2.Ancestor(height)
			if ancestor == nil {
				t.Fatalf("Ancestor(%d) returned nil", height)
			}
			if got := ancestor.Height(); got != height {
				t.Fatalf("Ancestor(%d) height = %d, want %d", height, got, height)
			}
		}

		for _, height := range []int32{-1, 3} {
			if ancestor := entry2.Ancestor(height); ancestor != nil {
				t.Fatalf("Ancestor(%d) returned height %d, want nil", height, ancestor.Height())
			}
		}
	})

	t.Run("GetHeader", func(t *testing.T) {
		// Get genesis block entry
		genesisEntry, err := chain.GetByHeight(0)
		if err != nil {
			t.Fatalf("GetByHeight(0) error = %v", err)
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

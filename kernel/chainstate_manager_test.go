package kernel

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestChainstateManagerCreation(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "bitcoin_kernel_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dataDir := filepath.Join(tempDir, "data")
	blocksDir := filepath.Join(tempDir, "blocks")

	// Create context
	ctx, err := NewDefaultContext()
	if err != nil {
		t.Fatalf("NewDefaultContext() error = %v", err)
	}
	defer ctx.Destroy()

	// Create chainstate manager options
	opts, err := NewChainstateManagerOptions(ctx, dataDir, blocksDir)
	if err != nil {
		t.Fatalf("NewChainstateManagerOptions() error = %v", err)
	}
	defer opts.Destroy()

	// Set some options
	opts.SetWorkerThreads(1)
	opts.SetBlockTreeDBInMemory(true)
	opts.SetChainstateDBInMemory(true)

	// Create chainstate manager
	manager, err := NewChainstateManager(ctx, opts)
	if err != nil {
		t.Fatalf("NewChainstateManager() error = %v", err)
	}
	defer manager.Destroy()

	// Test getting genesis block index
	genesisIndex, err := manager.GetBlockIndexFromGenesis()
	if err != nil {
		t.Fatalf("GetBlockIndexFromGenesis() error = %v", err)
	}
	defer genesisIndex.Destroy()

	// Test getting genesis block height
	height := genesisIndex.Height()
	if height != 0 {
		t.Errorf("Expected genesis height 0, got %d", height)
	}

	// Test getting genesis block hash
	genesisHash, err := genesisIndex.Hash()
	if err != nil {
		t.Fatalf("BlockIndex.Hash() error = %v", err)
	}
	defer genesisHash.Destroy()

	hashBytes := genesisHash.Bytes()
	if len(hashBytes) != 32 {
		t.Errorf("Expected hash length 32, got %d", len(hashBytes))
	}
}

// TestLoadAndValidateBlock demonstrates the complete workflow of:
// 1. Setting up a chainstate manager
// 2. Reading a block from disk (using genesis block index)
// 3. Validating a block
func TestLoadAndValidateBlock(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "bitcoin_kernel_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dataDir := filepath.Join(tempDir, "data")
	blocksDir := filepath.Join(tempDir, "blocks")

	// Create context
	ctx, err := NewDefaultContext()
	if err != nil {
		t.Fatalf("NewDefaultContext() error = %v", err)
	}
	defer ctx.Destroy()

	// Create chainstate manager options
	opts, err := NewChainstateManagerOptions(ctx, dataDir, blocksDir)
	if err != nil {
		t.Fatalf("NewChainstateManagerOptions() error = %v", err)
	}
	defer opts.Destroy()

	// Configure for in-memory operation
	opts.SetWorkerThreads(1)
	opts.SetBlockTreeDBInMemory(true)
	opts.SetChainstateDBInMemory(true)

	// Create chainstate manager
	manager, err := NewChainstateManager(ctx, opts)
	if err != nil {
		t.Fatalf("NewChainstateManager() error = %v", err)
	}
	defer manager.Destroy()

	// Initialize empty databases
	err = manager.ImportBlocks(nil)
	if err != nil {
		t.Fatalf("ImportBlocks() error = %v", err)
	}

	// Test 1: Get genesis block index
	t.Log("=== Test 1: Getting genesis block index ===")
	genesisIndex, err := manager.GetBlockIndexFromGenesis()
	if err != nil {
		t.Fatalf("GetBlockIndexFromGenesis() error = %v", err)
	}
	defer genesisIndex.Destroy()

	height := genesisIndex.Height()
	t.Logf("Genesis block height: %d", height)
	if height != 0 {
		t.Errorf("Expected genesis height 0, got %d", height)
	}

	// Test 2: Read genesis block from disk
	t.Log("=== Test 2: Reading genesis block from disk ===")
	genesisBlock, err := manager.ReadBlockFromDisk(genesisIndex)
	if err != nil {
		t.Fatalf("ReadBlockFromDisk() error = %v", err)
	}
	defer genesisBlock.Destroy()

	// Get block hash
	blockHash, err := genesisBlock.Hash()
	if err != nil {
		t.Fatalf("Block.Hash() error = %v", err)
	}
	defer blockHash.Destroy()
	t.Logf("Genesis block hash: %x", ReverseBytes(blockHash.Bytes()))

	// Test 3: Process/validate a new block
	t.Log("=== Test 3: Creating and validating a block ===")

	// Use the complete genesis block data for validation test
	genesisHex := "0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c0101000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000"
	genesisBytes, err := hex.DecodeString(genesisHex)
	if err != nil {
		t.Fatalf("Failed to decode genesis hex: %v", err)
	}

	testBlock, err := NewBlockFromRaw(genesisBytes)
	if err != nil {
		t.Fatalf("NewBlockFromRaw() error = %v", err)
	}
	defer testBlock.Destroy()

	// Process the block (validate it)
	success, isNewBlock, err := manager.ProcessBlock(testBlock)
	if err != nil {
		t.Fatalf("ProcessBlock() error = %v", err)
	}

	t.Logf("Block validation successful: %v", success)
	t.Logf("Was new block: %v", isNewBlock)

	// Test 4: Navigate block index
	t.Log("=== Test 4: Block index navigation ===")

	// Get current tip
	tipIndex, err := manager.GetBlockIndexFromTip()
	if err != nil {
		t.Fatalf("GetBlockIndexFromTip() error = %v", err)
	}
	defer tipIndex.Destroy()
	t.Logf("Current tip height: %d", tipIndex.Height())

	// Get block by height
	blockAt0, err := manager.GetBlockIndexFromHeight(0)
	if err != nil {
		t.Fatalf("GetBlockIndexFromHeight(0) error = %v", err)
	}
	defer blockAt0.Destroy()
	t.Logf("Block at height 0: %d", blockAt0.Height())

	// Test previous block navigation (should be nil for genesis)
	prevBlock := genesisIndex.Previous()
	if prevBlock != nil {
		defer prevBlock.Destroy()
		t.Error("Expected no previous block for genesis, but got one")
	}

	// Test next block navigation
	nextBlock, err := manager.GetNextBlockIndex(genesisIndex)
	if err != nil {
		t.Fatalf("GetNextBlockIndex() error = %v", err)
	}
	if nextBlock != nil {
		defer nextBlock.Destroy()
		t.Logf("Next block height: %d", nextBlock.Height())
	} else {
		t.Log("No next block (genesis is tip)")
	}

	t.Log("=== Load and validate test completed successfully ===")
}

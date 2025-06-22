package kernel

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
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

	ctx, err := NewDefaultContext()
	if err != nil {
		t.Fatalf("NewDefaultContext() error = %v", err)
	}
	defer ctx.Destroy()

	opts, err := NewChainstateManagerOptions(ctx, dataDir, blocksDir)
	if err != nil {
		t.Fatalf("NewChainstateManagerOptions() error = %v", err)
	}
	defer opts.Destroy()

	opts.SetWorkerThreads(1)
	opts.SetBlockTreeDBInMemory(true)
	opts.SetChainstateDBInMemory(true)

	manager, err := NewChainstateManager(ctx, opts)
	if err != nil {
		t.Fatalf("NewChainstateManager() error = %v", err)
	}
	defer manager.Destroy()

	genesisIndex, err := manager.GetBlockIndexFromGenesis()
	if err != nil {
		t.Fatalf("GetBlockIndexFromGenesis() error = %v", err)
	}
	defer genesisIndex.Destroy()

	height := genesisIndex.Height()
	if height != 0 {
		t.Errorf("Expected genesis height 0, got %d", height)
	}

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

// TestLoadAndValidateBlock demonstrates the following workflow:
// 1. Setting up a chainstate manager with regtest parameters
// 2. Loading blocks from data/regtest/blocks.txt
// 3. Validating and processing multiple blocks in sequence
func TestLoadAndValidateBlock(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "bitcoin_kernel_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dataDir := filepath.Join(tempDir, "data")
	blocksDir := filepath.Join(tempDir, "blocks")

	// Create context with regtest chain parameters
	contextOpts, err := NewContextOptions()
	if err != nil {
		t.Fatalf("NewContextOptions() error = %v", err)
	}
	defer contextOpts.Destroy()

	// Set regtest chain parameters
	chainParams, err := NewChainParameters(ChainTypeRegtest)
	if err != nil {
		t.Fatalf("NewChainParameters() error = %v", err)
	}
	defer chainParams.Destroy()

	contextOpts.SetChainParams(chainParams)

	ctx, err := NewContext(contextOpts)
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
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

	// Load block data from data/regtest/blocks.txt
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	projectRoot := filepath.Dir(wd)
	blocksFile := filepath.Join(projectRoot, "data", "regtest", "blocks.txt")

	blocksData, err := os.ReadFile(blocksFile)
	if err != nil {
		t.Fatalf("Failed to read blocks file: %v", err)
	}

	var blockLines []string
	for _, line := range strings.Split(string(blocksData), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			blockLines = append(blockLines, line)
		}
	}
	if len(blockLines) == 0 {
		t.Fatal("No block data found in blocks.txt")
	}
	t.Logf("Found %d blocks in regtest data", len(blockLines))

	// Process all blocks from the data file
	for i := 0; i < len(blockLines); i++ {
		blockHex := blockLines[i]

		// Decode hex data
		blockBytes, err := hex.DecodeString(blockHex)
		if err != nil {
			t.Fatalf("Failed to decode block %d hex: %v", i+1, err)
		}

		// Create block from raw data
		block, err := NewBlockFromRaw(blockBytes)
		if err != nil {
			t.Fatalf("NewBlockFromRaw() failed for block %d: %v", i+1, err)
		}
		defer block.Destroy()

		// Process the block (validate it)
		success, isNewBlock, err := manager.ProcessBlock(block)
		if err != nil {
			t.Fatalf("ProcessBlock() failed for block %d: %v", i+1, err)
		}
		if !success || !isNewBlock {
			t.Fatalf("ProcessBlock() failed for block %d", i+1)
		}

		// Assert tip height
		tipIndex, err := manager.GetBlockIndexFromTip()
		if err != nil {
			t.Fatalf("GetBlockIndexFromTip() error = %v", err)
		}
		defer tipIndex.Destroy()

		expectedTipHeight := int32(i) + 1
		if tipIndex.Height() != expectedTipHeight {
			t.Fatalf("Expected tip height %d; got %d", expectedTipHeight, tipIndex.Height())
		}

		// Assert tip hash
		hash, err := block.Hash()
		if err != nil {
			t.Fatalf("Failed to get last block hash: %v", err)
		}
		defer hash.Destroy()

		hashHex := hex.EncodeToString(hash.Bytes())

		tipHash, err := tipIndex.Hash()
		if err != nil {
			t.Fatalf("Failed to get tip hash: %v", err)
		}
		defer tipHash.Destroy()

		tipHashHex := hex.EncodeToString(tipHash.Bytes())

		if hashHex != tipHashHex {
			t.Fatalf("Expected tip hash %s; got %s", hashHex, tipHashHex)
		}
	}
}

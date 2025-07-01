package kernel

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChainstateManager(t *testing.T) {
	suite := SetupChainstateManagerTestSuite(t)

	t.Run("genesis validation", suite.TestGenesis)
	t.Run("tip validation", suite.TestTip)
}

func (s *ChainstateManagerTestSuite) TestGenesis(t *testing.T) {
	genesisIndex, err := s.Manager.GetBlockIndexFromGenesis()
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

func (s *ChainstateManagerTestSuite) TestTip(t *testing.T) {
	tipIndex, err := s.Manager.GetBlockIndexFromTip()
	if err != nil {
		t.Fatalf("GetBlockIndexFromTip() error = %v", err)
	}
	defer tipIndex.Destroy()

	height := tipIndex.Height()
	if height <= 0 {
		t.Errorf("Expected tip height > 0, got %d", height)
	}

	tipHash, err := tipIndex.Hash()
	if err != nil {
		t.Fatalf("Failed to get tip hash: %v", err)
	}
	defer tipHash.Destroy()

	hashBytes := tipHash.Bytes()
	if len(hashBytes) != 32 {
		t.Errorf("Expected hash length 32, got %d", len(hashBytes))
	}

	if tipIndex.Height() != s.ImportedBlocksCount {
		t.Errorf("Expected tip height %d, got %d", s.ImportedBlocksCount, tipIndex.Height())
	}
}

type ChainstateManagerTestSuite struct {
	Manager             *ChainstateManager
	ImportedBlocksCount int32
}

func SetupChainstateManagerTestSuite(t *testing.T) *ChainstateManagerTestSuite {
	tempDir, err := os.MkdirTemp("", "bitcoin_kernel_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	dataDir := filepath.Join(tempDir, "data")
	blocksDir := filepath.Join(tempDir, "blocks")

	contextOpts, err := NewContextOptions()
	if err != nil {
		t.Fatalf("NewContextOptions() error = %v", err)
	}
	t.Cleanup(func() { contextOpts.Destroy() })

	chainParams, err := NewChainParameters(ChainTypeRegtest)
	if err != nil {
		t.Fatalf("NewChainParameters() error = %v", err)
	}
	t.Cleanup(func() { chainParams.Destroy() })

	contextOpts.SetChainParams(chainParams)

	ctx, err := NewContext(contextOpts)
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	t.Cleanup(func() { ctx.Destroy() })

	opts, err := NewChainstateManagerOptions(ctx, dataDir, blocksDir)
	if err != nil {
		t.Fatalf("NewChainstateManagerOptions() error = %v", err)
	}
	t.Cleanup(func() { opts.Destroy() })

	opts.SetWorkerThreads(1)
	opts.SetBlockTreeDBInMemory(true)
	opts.SetChainstateDBInMemory(true)

	// Create chainstate manager
	manager, err := NewChainstateManager(ctx, opts)
	if err != nil {
		t.Fatalf("NewChainstateManager() error = %v", err)
	}
	t.Cleanup(func() { manager.Destroy() })

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

	for i := 0; i < len(blockLines); i++ {
		blockHex := blockLines[i]

		blockBytes, err := hex.DecodeString(blockHex)
		if err != nil {
			t.Fatalf("Failed to decode block %d hex: %v", i+1, err)
		}

		block, err := NewBlockFromRaw(blockBytes)
		if err != nil {
			t.Fatalf("NewBlockFromRaw() failed for block %d: %v", i+1, err)
		}
		defer block.Destroy()

		success, isNewBlock, err := manager.ProcessBlock(block)
		if err != nil {
			t.Fatalf("ProcessBlock() failed for block %d: %v", i+1, err)
		}
		if !success || !isNewBlock {
			t.Fatalf("ProcessBlock() failed for block %d", i+1)
		}
	}

	return &ChainstateManagerTestSuite{
		Manager:             manager,
		ImportedBlocksCount: int32(len(blockLines)),
	}
}

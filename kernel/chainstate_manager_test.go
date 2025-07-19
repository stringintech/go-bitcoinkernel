package kernel

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChainstateManager(t *testing.T) {
	suite := ChainstateManagerTestSuite{
		MaxBlockHeightToImport: 0,   // load all blocks from data/regtest/block.txt
		NotificationCallbacks:  nil, // no notification callbacks
		ValidationCallbacks:    nil, // no validation callbacks
	}
	suite.Setup(t)

	t.Run("genesis validation", suite.TestGenesis)
	t.Run("tip validation", suite.TestTip)
	t.Run("block undo", suite.TestBlockUndo)
}

func (s *ChainstateManagerTestSuite) TestGenesis(t *testing.T) {
	genesisIndex, err := s.Manager.GetBlockIndexGenesis()
	if err != nil {
		t.Fatalf("GetBlockIndexGenesis() error = %v", err)
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
	tipIndex, err := s.Manager.GetBlockIndexTip()
	if err != nil {
		t.Fatalf("GetBlockIndexTip() error = %v", err)
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

func (s *ChainstateManagerTestSuite) TestBlockUndo(t *testing.T) {
	blockIndex, err := s.Manager.GetBlockIndexByHeight(202)
	if err != nil {
		t.Fatalf("GetBlockIndexByHeight(202) error = %v", err)
	}
	defer blockIndex.Destroy()

	blockUndo, err := s.Manager.ReadBlockUndo(blockIndex)
	if err != nil {
		t.Fatalf("ReadBlockUndo() error = %v", err)
	}
	defer blockUndo.Destroy()

	// Test transaction count
	txCount := blockUndo.Size()
	if txCount != 20 {
		t.Errorf("Expected 20 transactions, got %d", txCount)
	}

	// Verify each transaction is a valid TransactionUndo
	for i := uint64(0); i < txCount; i++ {
		undoSize := blockUndo.GetTransactionUndoSize(i)
		if undoSize != 1 {
			t.Errorf("Expected transaction undo size 1, got %d", undoSize)
		}

		_, err := blockUndo.GetUndoOutputByIndex(i, 0)
		if err != nil {
			t.Fatalf("GetUndoOutputByIndex() error = %v", err)
		}

		height := blockUndo.GetUndoOutputHeightByIndex(i, 0)
		if height <= 0 {
			t.Fatalf("GetUndoOutputHeightByIndex() height %d, want > 0", height)
		}
	}
}

type ChainstateManagerTestSuite struct {
	MaxBlockHeightToImport int32 // leave zero to load all blocks
	NotificationCallbacks  *NotificationCallbacks
	ValidationCallbacks    *ValidationInterfaceCallbacks

	Manager             *ChainstateManager
	ImportedBlocksCount int32
}

func (s *ChainstateManagerTestSuite) Setup(t *testing.T) {
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

	if s.NotificationCallbacks != nil {
		err = contextOpts.SetNotifications(s.NotificationCallbacks)
		if err != nil {
			t.Fatalf("SetNotifications() error = %v", err)
		}
	}

	if s.ValidationCallbacks != nil {
		err = contextOpts.SetValidationInterface(s.ValidationCallbacks)
		if err != nil {
			t.Fatalf("SetValidationInterface() error = %v", err)
		}
	}

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
		if s.MaxBlockHeightToImport != 0 && len(blockLines) >= int(s.MaxBlockHeightToImport) {
			break
		}
	}
	if len(blockLines) == 0 {
		t.Fatal("No block data found in blocks.txt")
	}

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

	s.Manager = manager
	s.ImportedBlocksCount = int32(len(blockLines))
}

package kernel

import (
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"slices"
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

	t.Run("read block", suite.TestReadBlock)
	t.Run("block undo", suite.TestBlockSpentOutputs)
	t.Run("transaction spent outputs", suite.TestTransactionSpentOutputs)
	t.Run("get block tree entry by hash", suite.TestGetBlockTreeEntryByHash)
}

func (s *ChainstateManagerTestSuite) TestBlockSpentOutputs(t *testing.T) {
	chain := s.Manager.GetActiveChain()
	blockIndex := chain.GetByHeight(202)

	blockSpentOutputs, err := s.Manager.ReadBlockSpentOutputs(blockIndex)
	if err != nil {
		t.Fatalf("ReadBlockSpentOutputs() error = %v", err)
	}
	defer blockSpentOutputs.Destroy()

	t.Run("Count", func(t *testing.T) {
		txCount := blockSpentOutputs.Count()
		if txCount != 20 {
			t.Errorf("Expected 20 transactions, got %d", txCount)
		}
	})

	t.Run("GetTransactionSpentOutputsAt", func(t *testing.T) {
		for i := uint64(0); i < blockSpentOutputs.Count(); i++ {
			_, err := blockSpentOutputs.GetTransactionSpentOutputsAt(i)
			if err != nil {
				t.Fatalf("GetTransactionSpentOutputsAt(%d) error = %v", i, err)
			}
		}
	})

	t.Run("TransactionsSpentOutputs", func(t *testing.T) {
		count := len(slices.Collect(blockSpentOutputs.TransactionsSpentOutputs()))
		if count != 20 {
			t.Errorf("Expected to iterate over 20 transaction spent outputs, got %d", count)
		}
	})

	t.Run("TransactionsSpentOutputsRange", func(t *testing.T) {
		count := len(slices.Collect(blockSpentOutputs.TransactionsSpentOutputsRange(0, 1000)))
		if count != 20 {
			t.Errorf("Expected to iterate over 20 transaction spent outputs, got %d", count)
		}

		count = len(slices.Collect(blockSpentOutputs.TransactionsSpentOutputsRange(10, 15)))
		if count != 5 {
			t.Errorf("Expected to iterate over 5 transaction spent outputs, got %d", count)
		}
	})

	t.Run("TransactionsSpentOutputsFrom", func(t *testing.T) {
		count := len(slices.Collect(blockSpentOutputs.TransactionsSpentOutputsFrom(0)))
		if count != 20 {
			t.Errorf("Expected to iterate over 20 transaction spent outputs, got %d", count)
		}

		count = len(slices.Collect(blockSpentOutputs.TransactionsSpentOutputsFrom(15)))
		if count != 5 {
			t.Errorf("Expected to iterate over 5 transaction spent outputs, got %d", count)
		}
	})
}

func (s *ChainstateManagerTestSuite) TestTransactionSpentOutputs(t *testing.T) {
	chain := s.Manager.GetActiveChain()
	blockIndex := chain.GetByHeight(202)

	blockSpentOutputs, err := s.Manager.ReadBlockSpentOutputs(blockIndex)
	if err != nil {
		t.Fatalf("ReadBlockSpentOutputs() error = %v", err)
	}
	defer blockSpentOutputs.Destroy()

	txSpentOutputs, err := blockSpentOutputs.GetTransactionSpentOutputsAt(0)
	if err != nil {
		t.Fatalf("GetTransactionSpentOutputsAt(0) error = %v", err)
	}

	t.Run("Count", func(t *testing.T) {
		if txSpentOutputs.Count() != 1 {
			t.Errorf("Expected 1, got %d", txSpentOutputs.Count())
		}
	})

	t.Run("GetCoinAt", func(t *testing.T) {
		coin, err := txSpentOutputs.GetCoinAt(0)
		if err != nil {
			t.Fatalf("GetCoinAt(0) error = %v", err)
		}

		coin.GetOutput()

		height := coin.ConfirmationHeight()
		if height <= 0 {
			t.Fatalf("ConfirmationHeight() height %d, want > 0", height)
		}

		_, err = txSpentOutputs.GetCoinAt(txSpentOutputs.Count())
		if !errors.Is(err, ErrKernelIndexOutOfBounds) {
			t.Errorf("Expected ErrKernelIndexOutOfBounds for out of bounds coin, got %v", err)
		}
	})

	t.Run("Coins", func(t *testing.T) {
		count := len(slices.Collect(txSpentOutputs.Coins()))
		if count != 1 {
			t.Errorf("Expected to iterate over 1 coin, got %d", count)
		}
	})

	t.Run("CoinsRange", func(t *testing.T) {
		count := len(slices.Collect(txSpentOutputs.CoinsRange(0, 1000)))
		if count != 1 {
			t.Errorf("Expected to iterate over 1 coin, got %d", count)
		}

		count = len(slices.Collect(txSpentOutputs.CoinsRange(1, 2)))
		if count != 0 {
			t.Errorf("Expected to iterate over 0 coins, got %d", count)
		}
	})

	t.Run("CoinsFrom", func(t *testing.T) {
		count := len(slices.Collect(txSpentOutputs.CoinsFrom(0)))
		if count != 1 {
			t.Errorf("Expected to iterate over 1 coin, got %d", count)
		}

		count = len(slices.Collect(txSpentOutputs.CoinsFrom(1)))
		if count != 0 {
			t.Errorf("Expected to iterate over 0 coins, got %d", count)
		}
	})
}

func (s *ChainstateManagerTestSuite) TestReadBlock(t *testing.T) {
	chain := s.Manager.GetActiveChain()

	// Test reading genesis block
	genesis := chain.GetByHeight(0)
	genesisBlock, err := s.Manager.ReadBlock(genesis)
	if err != nil {
		t.Fatalf("ChainstateManager.ReadBlock() for genesis error = %v", err)
	}
	if genesisBlock == nil {
		t.Fatal("Read genesis block is nil")
	}
	defer genesisBlock.Destroy()

	// Test reading tip block
	tip := chain.GetByHeight(chain.GetHeight())
	tipBlock, err := s.Manager.ReadBlock(tip)
	if err != nil {
		t.Fatalf("ChainstateManager.ReadBlock() for tip error = %v", err)
	}
	if tipBlock == nil {
		t.Fatal("Read tip block is nil")
	}
	defer tipBlock.Destroy()
}

func (s *ChainstateManagerTestSuite) TestGetBlockTreeEntryByHash(t *testing.T) {
	chain := s.Manager.GetActiveChain()

	// Test getting genesis block by hash
	genesis := chain.GetByHeight(0)

	genesisHash := genesis.Hash()

	// Use GetBlockTreeEntryByHash to find genesis
	foundGenesisIndex := s.Manager.GetBlockTreeEntryByHash(genesisHash)
	if foundGenesisIndex == nil {
		t.Fatal("Found genesis block tree entry is nil")
	}

	// Verify found block has same height as original
	foundHeight := foundGenesisIndex.Height()
	originalHeight := genesis.Height()
	if foundHeight != originalHeight {
		t.Errorf("Found genesis height %d, expected %d", foundHeight, originalHeight)
	}

	// Test getting tip block by hash
	tipIndex := chain.GetByHeight(chain.GetHeight())

	tipHash := tipIndex.Hash()

	foundTipIndex := s.Manager.GetBlockTreeEntryByHash(tipHash)
	if foundTipIndex == nil {
		t.Fatal("Found tip block tree entry is nil")
	}

	// Verify found tip has same height as original
	foundTipHeight := foundTipIndex.Height()
	originalTipHeight := tipIndex.Height()
	if foundTipHeight != originalTipHeight {
		t.Errorf("Found tip height %d, expected %d", foundTipHeight, originalTipHeight)
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
	t.Cleanup(func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("Failed to remove temp dir: %v", err)
		}
	})

	dataDir := filepath.Join(tempDir, "data")
	blocksDir := filepath.Join(tempDir, "blocks")

	var contextOpts []ContextOption
	contextOpts = append(contextOpts, WithChainType(ChainTypeRegtest))

	if s.NotificationCallbacks != nil {
		contextOpts = append(contextOpts, WithNotifications(s.NotificationCallbacks))
	}

	if s.ValidationCallbacks != nil {
		contextOpts = append(contextOpts, WithValidationInterface(s.ValidationCallbacks))
	}

	ctx, err := NewContext(contextOpts...)
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	t.Cleanup(func() { ctx.Destroy() })

	manager, err := NewChainstateManager(ctx, dataDir, blocksDir,
		WithWorkerThreads(1),
		WithBlockTreeDBInMemory(true),
		WithChainstateDBInMemory(),
		WithWipeDBs(true, true),
	)
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

		block, err := NewBlock(blockBytes)
		if err != nil {
			t.Fatalf("NewBlockFromRaw() failed for block %d: %v", i+1, err)
		}
		defer block.Destroy()

		ok, duplicate := manager.ProcessBlock(block)
		if !ok || duplicate {
			t.Fatalf("ProcessBlock() failed for block %d", i+1)
		}
	}

	s.Manager = manager
	s.ImportedBlocksCount = int32(len(blockLines))
}

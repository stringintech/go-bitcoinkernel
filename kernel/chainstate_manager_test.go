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

	t.Run("read block", suite.TestReadBlock)
	t.Run("block undo", suite.TestBlockSpentOutputs)
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

	// Test transaction spent outputs count
	txCount := blockSpentOutputs.Count()
	if txCount != 20 {
		t.Errorf("Expected 20 transactions, got %d", txCount)
	}

	// Verify each transaction spent outputs
	for i := uint64(0); i < txCount; i++ {
		txSpentOutputs, err := blockSpentOutputs.GetTransactionSpentOutputsAt(i)
		if err != nil {
			t.Fatalf("GetTransactionSpentOutputsAt(%d) error = %v", i, err)
		}

		spentOutputSize := txSpentOutputs.Count()
		if spentOutputSize != 1 {
			t.Errorf("Expected transaction spent output size 1, got %d", spentOutputSize)
		}

		coin, err := txSpentOutputs.GetCoinAt(0)
		if err != nil {
			t.Fatalf("GetCoinAt(0) error = %v", err)
		}

		coin.GetOutput()

		height := coin.ConfirmationHeight()
		if height <= 0 {
			t.Fatalf("ConfirmationHeight() height %d, want > 0", height)
		}
	}
}

func (s *ChainstateManagerTestSuite) TestReadBlock(t *testing.T) {
	chain := s.Manager.GetActiveChain()

	// Test reading genesis block
	genesis := chain.GetGenesis()
	genesisBlock, err := s.Manager.ReadBlock(genesis)
	if err != nil {
		t.Fatalf("ChainstateManager.ReadBlock() for genesis error = %v", err)
	}
	if genesisBlock == nil {
		t.Fatal("Read genesis block is nil")
	}
	defer genesisBlock.Destroy()

	// Test reading tip block
	tip := chain.GetTip()
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
	genesis := chain.GetGenesis()

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
	tipIndex := chain.GetTip()

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
		WithBlockTreeDBInMemory,
		WithChainstateDBInMemory,
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

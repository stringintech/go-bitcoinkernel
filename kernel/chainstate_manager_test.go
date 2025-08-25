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
	chain, err := s.Manager.GetActiveChain()
	if err != nil {
		t.Fatalf("GetActiveChain() error = %v", err)
	}
	defer chain.Destroy()

	blockIndex := chain.GetByHeight(202)
	defer blockIndex.Destroy()

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
		defer txSpentOutputs.Destroy()

		spentOutputSize := txSpentOutputs.Count()
		if spentOutputSize != 1 {
			t.Errorf("Expected transaction spent output size 1, got %d", spentOutputSize)
		}

		coin, err := txSpentOutputs.GetCoinAt(0)
		if err != nil {
			t.Fatalf("GetCoinAt(0) error = %v", err)
		}
		defer coin.Destroy()

		_, err = coin.GetOutput()
		if err != nil {
			t.Fatalf("GetOutput() error = %v", err)
		}

		height := coin.ConfirmationHeight()
		if height <= 0 {
			t.Fatalf("ConfirmationHeight() height %d, want > 0", height)
		}
	}
}

func (s *ChainstateManagerTestSuite) TestReadBlock(t *testing.T) {
	chain, err := s.Manager.GetActiveChain()
	if err != nil {
		t.Fatalf("GetActiveChain() error = %v", err)
	}
	defer chain.Destroy()

	// Test reading genesis block
	genesis, err := chain.GetGenesis()
	if err != nil {
		t.Fatalf("GetGenesis() error = %v", err)
	}
	defer genesis.Destroy()

	genesisBlock, err := s.Manager.ReadBlock(genesis)
	if err != nil {
		t.Fatalf("ChainstateManager.ReadBlock() for genesis error = %v", err)
	}
	if genesisBlock == nil {
		t.Fatal("Read genesis block is nil")
	}
	defer genesisBlock.Destroy()

	// Verify genesis block has expected properties
	genesisHash, err := genesisBlock.Hash()
	if err != nil {
		t.Fatalf("Genesis block Hash() error = %v", err)
	}
	defer genesisHash.Destroy()

	hashBytes := genesisHash.Bytes()
	if len(hashBytes) != 32 {
		t.Errorf("Expected hash length 32, got %d", len(hashBytes))
	}

	// Test reading tip block
	tip, err := chain.GetTip()
	if err != nil {
		t.Fatalf("GetTip() error = %v", err)
	}
	defer tip.Destroy()

	tipBlock, err := s.Manager.ReadBlock(tip)
	if err != nil {
		t.Fatalf("ChainstateManager.ReadBlock() for tip error = %v", err)
	}
	if tipBlock == nil {
		t.Fatal("Read tip block is nil")
	}
	defer tipBlock.Destroy()

	// Verify tip block properties
	tipHash, err := tipBlock.Hash()
	if err != nil {
		t.Fatalf("Tip block Hash() error = %v", err)
	}
	defer tipHash.Destroy()

	tipHashBytes := tipHash.Bytes()
	if len(tipHashBytes) != 32 {
		t.Errorf("Expected tip hash length 32, got %d", len(tipHashBytes))
	}
}

func (s *ChainstateManagerTestSuite) TestGetBlockTreeEntryByHash(t *testing.T) {
	chain, err := s.Manager.GetActiveChain()
	if err != nil {
		t.Fatalf("GetActiveChain() error = %v", err)
	}
	defer chain.Destroy()

	// Test getting genesis block by hash
	genesis, err := chain.GetGenesis()
	if err != nil {
		t.Fatalf("GetGenesis() error = %v", err)
	}
	defer genesis.Destroy()

	genesisHash, err := genesis.Hash()
	if err != nil {
		t.Fatalf("Genesis Hash() error = %v", err)
	}
	defer genesisHash.Destroy()

	// Use GetBlockTreeEntryByHash to find genesis
	foundGenesisIndex, err := s.Manager.GetBlockTreeEntryByHash(genesisHash)
	if err != nil {
		t.Fatalf("ChainstateManager.GetBlockTreeEntryByHash() for genesis error = %v", err)
	}
	if foundGenesisIndex == nil {
		t.Fatal("Found genesis block tree entry is nil")
	}
	defer foundGenesisIndex.Destroy()

	// Verify found block has same height as original
	foundHeight := foundGenesisIndex.Height()
	originalHeight := genesis.Height()
	if foundHeight != originalHeight {
		t.Errorf("Found genesis height %d, expected %d", foundHeight, originalHeight)
	}

	// Test getting tip block by hash
	tipIndex, err := chain.GetTip()
	if err != nil {
		t.Fatalf("GetTip() error = %v", err)
	}
	defer tipIndex.Destroy()

	tipHash, err := tipIndex.Hash()
	if err != nil {
		t.Fatalf("Tip Hash() error = %v", err)
	}
	defer tipHash.Destroy()

	foundTipIndex, err := s.Manager.GetBlockTreeEntryByHash(tipHash)
	if err != nil {
		t.Fatalf("ChainstateManager.GetBlockTreeEntryByHash() for tip error = %v", err)
	}
	if foundTipIndex == nil {
		t.Fatal("Found tip block tree entry is nil")
	}
	defer foundTipIndex.Destroy()

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
	// Wipe both databases to enable proper initialization
	opts.SetWipeDBs(true, true)

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

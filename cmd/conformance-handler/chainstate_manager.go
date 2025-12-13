package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// handleChainstateManagerCreate creates a chainstate manager from a context
func handleChainstateManagerCreate(registry *Registry, req Request) Response {
	var params struct {
		Context RefObject `json:"context"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsResponse(req.ID, "failed to parse params")
	}

	if req.Ref == "" {
		return NewInvalidParamsResponse(req.ID, "ref field is required")
	}

	// Get context from registry
	ctx, err := registry.GetContext(params.Context.Ref)
	if err != nil {
		return NewInvalidParamsResponse(req.ID, err.Error())
	}

	// Create temp directory for chainstate data
	tempDir, err := os.MkdirTemp("", "btck_conformance_test_*")
	if err != nil {
		return NewEmptyErrorResponse(req.ID)
	}

	dataDir := filepath.Join(tempDir, "data")
	blocksDir := filepath.Join(tempDir, "blocks")

	// Create chainstate manager
	manager, err := kernel.NewChainstateManager(ctx, dataDir, blocksDir)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		return NewEmptyErrorResponse(req.ID)
	}

	registry.Store(req.Ref, &ChainstateManagerState{
		Manager: manager,
		TempDir: tempDir,
	})

	return NewSuccessResponseWithRef(req.ID, req.Ref)
}

// handleChainstateManagerGetActiveChain gets the active chain from a chainstate manager
func handleChainstateManagerGetActiveChain(registry *Registry, req Request) Response {
	var params struct {
		ChainstateManager RefObject `json:"chainstate_manager"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsResponse(req.ID, "failed to parse params")
	}

	if req.Ref == "" {
		return NewInvalidParamsResponse(req.ID, "ref field is required")
	}

	// Get chainstate manager from registry
	csm, err := registry.GetChainstateManager(params.ChainstateManager.Ref)
	if err != nil {
		return NewInvalidParamsResponse(req.ID, err.Error())
	}

	// Get active chain
	chain := csm.Manager.GetActiveChain()

	registry.Store(req.Ref, chain)

	return NewSuccessResponseWithRef(req.ID, req.Ref)
}

// handleChainstateManagerProcessBlock processes a block
func handleChainstateManagerProcessBlock(registry *Registry, req Request) Response {
	var params struct {
		ChainstateManager RefObject `json:"chainstate_manager"`
		Block             RefObject `json:"block"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsResponse(req.ID, "failed to parse params")
	}

	// Get chainstate manager from registry
	csm, err := registry.GetChainstateManager(params.ChainstateManager.Ref)
	if err != nil {
		return NewInvalidParamsResponse(req.ID, err.Error())
	}

	// Get block from registry
	block, err := registry.GetBlock(params.Block.Ref)
	if err != nil {
		return NewInvalidParamsResponse(req.ID, err.Error())
	}

	// Process the block
	ok, newBlock := csm.Manager.ProcessBlock(block)
	if !ok {
		return NewEmptyErrorResponse(req.ID)
	}

	// Return result with new_block field
	result := struct {
		NewBlock bool `json:"new_block"`
	}{
		NewBlock: newBlock,
	}
	return NewSuccessResponse(req.ID, result)
}

// handleChainstateManagerDestroy destroys a chainstate manager
func handleChainstateManagerDestroy(registry *Registry, req Request) Response {
	var params struct {
		ChainstateManager RefObject `json:"chainstate_manager"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsResponse(req.ID, "failed to parse params")
	}

	// Destroy and remove from registry
	if err := registry.Destroy(params.ChainstateManager.Ref); err != nil {
		return NewInvalidParamsResponse(req.ID, err.Error())
	}

	return NewEmptySuccessResponse(req.ID)
}

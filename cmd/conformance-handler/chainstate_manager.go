package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// handleChainstateManagerCreate creates a chainstate manager from a context
func handleChainstateManagerCreate(registry *Registry, req Request) (Response, error) {
	var params struct {
		Context RefObject `json:"context"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	ctx, err := registry.GetContext(params.Context.Ref)
	if err != nil {
		return Response{}, err
	}

	tempDir, err := os.MkdirTemp("", "btck_conformance_test_*")
	if err != nil {
		return NewEmptyErrorResponse(), nil
	}

	dataDir := filepath.Join(tempDir, "data")
	blocksDir := filepath.Join(tempDir, "blocks")

	manager, err := kernel.NewChainstateManager(ctx, dataDir, blocksDir)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		return NewEmptyErrorResponse(), nil
	}

	registry.Store(req.Ref, &ChainstateManagerState{
		Manager: manager,
		TempDir: tempDir,
	})

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleChainstateManagerGetActiveChain gets the active chain from a chainstate manager
func handleChainstateManagerGetActiveChain(registry *Registry, req Request) (Response, error) {
	var params struct {
		ChainstateManager RefObject `json:"chainstate_manager"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	csm, err := registry.GetChainstateManager(params.ChainstateManager.Ref)
	if err != nil {
		return Response{}, err
	}

	chain := csm.Manager.GetActiveChain()

	registry.Store(req.Ref, chain)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleChainstateManagerProcessBlock processes a block
func handleChainstateManagerProcessBlock(registry *Registry, req Request) (Response, error) {
	var params struct {
		ChainstateManager RefObject `json:"chainstate_manager"`
		Block             RefObject `json:"block"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	csm, err := registry.GetChainstateManager(params.ChainstateManager.Ref)
	if err != nil {
		return Response{}, err
	}

	block, err := registry.GetBlock(params.Block.Ref)
	if err != nil {
		return Response{}, err
	}

	ok, newBlock := csm.Manager.ProcessBlock(block)
	if !ok {
		return NewEmptyErrorResponse(), nil
	}

	result := struct {
		NewBlock bool `json:"new_block"`
	}{
		NewBlock: newBlock,
	}
	return NewSuccessResponse(result), nil
}

// handleChainstateManagerDestroy destroys a chainstate manager
func handleChainstateManagerDestroy(registry *Registry, req Request) (Response, error) {
	var params struct {
		ChainstateManager RefObject `json:"chainstate_manager"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if err := registry.Destroy(params.ChainstateManager.Ref); err != nil {
		return Response{}, err
	}

	return NewEmptySuccessResponse(), nil
}

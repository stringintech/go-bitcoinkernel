package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// handleBlockHashCreate creates a BlockHash from raw stored 32 bytes encoded as hex.
func handleBlockHashCreate(registry *Registry, req Request) (Response, error) {
	var params struct {
		BlockHash string `json:"block_hash"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	var hashBytes [32]byte
	if len(params.BlockHash) != 64 {
		return Response{}, fmt.Errorf("block_hash must be exactly 64 hex characters")
	}
	if _, err := hex.Decode(hashBytes[:], []byte(params.BlockHash)); err != nil {
		return Response{}, fmt.Errorf("block_hash must be valid hex: %w", err)
	}

	bh := kernel.NewBlockHash(hashBytes)
	registry.Store(req.Ref, bh)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleBlockHashToBytes returns the raw stored 32-byte block hash as a hex string.
func handleBlockHashToBytes(registry *Registry, req Request) (Response, error) {
	var params struct {
		BlockHash RefObject `json:"block_hash"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	bh, err := registry.GetBlockHash(params.BlockHash.Ref)
	if err != nil {
		return Response{}, err
	}

	raw := bh.Bytes()
	return NewSuccessResponse(fmt.Sprintf("%x", raw)), nil
}

// handleBlockHashEquals checks if two block hashes are equal
func handleBlockHashEquals(registry *Registry, req Request) (Response, error) {
	var params struct {
		Hash1 RefObject `json:"hash1"`
		Hash2 RefObject `json:"hash2"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	bh, err := registry.GetBlockHash(params.Hash1.Ref)
	if err != nil {
		return Response{}, err
	}

	bh2, err := registry.GetBlockHash(params.Hash2.Ref)
	if err != nil {
		return Response{}, err
	}

	return NewSuccessResponse(bh.Equals(bh2)), nil
}

// handleBlockHashCopy copies a block hash and stores the copy in the registry
func handleBlockHashCopy(registry *Registry, req Request) (Response, error) {
	var params struct {
		BlockHash RefObject `json:"block_hash"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	bh, err := registry.GetBlockHash(params.BlockHash.Ref)
	if err != nil {
		return Response{}, err
	}

	registry.Store(req.Ref, bh.Copy())

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleBlockHashDestroy destroys a block hash from the registry
func handleBlockHashDestroy(registry *Registry, req Request) (Response, error) {
	var params struct {
		BlockHash RefObject `json:"block_hash"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if _, err := registry.GetBlockHash(params.BlockHash.Ref); err != nil {
		return Response{}, err
	}

	if err := registry.Destroy(params.BlockHash.Ref); err != nil {
		return Response{}, err
	}

	return NewEmptySuccessResponse(), nil
}

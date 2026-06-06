package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// handleBlockHeaderCreate creates a block header from raw hex bytes
func handleBlockHeaderCreate(registry *Registry, req Request) (Response, error) {
	var params struct {
		RawBlockHeader string `json:"raw_block_header"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	headerBytes, err := hex.DecodeString(params.RawBlockHeader)
	if err != nil {
		return Response{}, fmt.Errorf("raw_block_header must be valid hex: %w", err)
	}

	header, err := kernel.NewBlockHeader(headerBytes)
	if err != nil {
		return NewEmptyErrorResponse(), nil
	}

	registry.Store(req.Ref, header)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleBlockHeaderGetHash gets the block hash from a block header and stores it in the registry
func handleBlockHeaderGetHash(registry *Registry, req Request) (Response, error) {
	var params struct {
		Header RefObject `json:"header"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	header, err := registry.GetBlockHeader(params.Header.Ref)
	if err != nil {
		return Response{}, err
	}

	hash := header.Hash()
	registry.Store(req.Ref, hash)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleBlockHeaderGetPrevHash gets the previous block hash view from a block header and stores it in the registry
func handleBlockHeaderGetPrevHash(registry *Registry, req Request) (Response, error) {
	var params struct {
		Header RefObject `json:"header"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	header, err := registry.GetBlockHeader(params.Header.Ref)
	if err != nil {
		return Response{}, err
	}

	prevHash := header.PrevHash()
	registry.Store(req.Ref, prevHash)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleBlockHeaderGetTimestamp returns the timestamp of a block header
func handleBlockHeaderGetTimestamp(registry *Registry, req Request) (Response, error) {
	var params struct {
		Header RefObject `json:"header"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	header, err := registry.GetBlockHeader(params.Header.Ref)
	if err != nil {
		return Response{}, err
	}

	return NewSuccessResponse(header.Timestamp()), nil
}

// handleBlockHeaderGetBits returns the nBits difficulty target of a block header
func handleBlockHeaderGetBits(registry *Registry, req Request) (Response, error) {
	var params struct {
		Header RefObject `json:"header"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	header, err := registry.GetBlockHeader(params.Header.Ref)
	if err != nil {
		return Response{}, err
	}

	return NewSuccessResponse(header.Bits()), nil
}

// handleBlockHeaderGetVersion returns the version of a block header
func handleBlockHeaderGetVersion(registry *Registry, req Request) (Response, error) {
	var params struct {
		Header RefObject `json:"header"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	header, err := registry.GetBlockHeader(params.Header.Ref)
	if err != nil {
		return Response{}, err
	}

	return NewSuccessResponse(header.Version()), nil
}

// handleBlockHeaderGetNonce returns the nonce of a block header
func handleBlockHeaderGetNonce(registry *Registry, req Request) (Response, error) {
	var params struct {
		Header RefObject `json:"header"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	header, err := registry.GetBlockHeader(params.Header.Ref)
	if err != nil {
		return Response{}, err
	}

	return NewSuccessResponse(header.Nonce()), nil
}

// handleBlockHeaderCopy copies a block header and stores the copy in the registry
func handleBlockHeaderCopy(registry *Registry, req Request) (Response, error) {
	var params struct {
		Header RefObject `json:"header"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	header, err := registry.GetBlockHeader(params.Header.Ref)
	if err != nil {
		return Response{}, err
	}

	headerCopy := header.Copy()
	registry.Store(req.Ref, headerCopy)

	return NewSuccessResponseWithRef(req.Ref), nil
}

func handleBlockHeaderToBytes(registry *Registry, req Request) (Response, error) {
	var params struct {
		Header RefObject `json:"header"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	header, err := registry.GetBlockHeader(params.Header.Ref)
	if err != nil {
		return Response{}, err
	}

	data, err := header.Bytes()
	if err != nil {
		return NewEmptyErrorResponse(), nil
	}

	return NewSuccessResponse(hex.EncodeToString(data[:])), nil
}

// handleBlockHeaderDestroy destroys a block header from the registry
func handleBlockHeaderDestroy(registry *Registry, req Request) (Response, error) {
	var params struct {
		Header RefObject `json:"header"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if err := registry.Destroy(params.Header.Ref); err != nil {
		return Response{}, err
	}

	return NewEmptySuccessResponse(), nil
}

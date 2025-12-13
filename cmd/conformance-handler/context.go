package main

import (
	"encoding/json"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// handleContextCreate creates a context with specified chain parameters
func handleContextCreate(registry *Registry, req Request) Response {
	var params struct {
		ChainParameters struct {
			ChainType string `json:"chain_type"`
		} `json:"chain_parameters"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsResponse(req.ID, "failed to parse params")
	}

	if req.Ref == "" {
		return NewInvalidParamsResponse(req.ID, "ref field is required")
	}

	// Parse chain type
	var chainType kernel.ChainType
	switch params.ChainParameters.ChainType {
	case "btck_ChainType_MAINNET":
		chainType = kernel.ChainTypeMainnet
	case "btck_ChainType_TESTNET":
		chainType = kernel.ChainTypeTestnet
	case "btck_ChainType_TESTNET_4":
		chainType = kernel.ChainTypeTestnet4
	case "btck_ChainType_SIGNET":
		chainType = kernel.ChainTypeSignet
	case "btck_ChainType_REGTEST":
		chainType = kernel.ChainTypeRegtest
	default:
		return NewInvalidParamsResponse(req.ID, "unknown chain_type: "+params.ChainParameters.ChainType)
	}

	// Create context
	ctx, err := kernel.NewContext(kernel.WithChainType(chainType))
	if err != nil {
		return NewEmptyErrorResponse(req.ID)
	}

	registry.Store(req.Ref, ctx)

	return NewSuccessResponseWithRef(req.ID, req.Ref)
}

// handleContextDestroy destroys a context
func handleContextDestroy(registry *Registry, req Request) Response {
	var params struct {
		Context RefObject `json:"context"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsResponse(req.ID, "failed to parse params")
	}

	// Destroy and remove from registry
	if err := registry.Destroy(params.Context.Ref); err != nil {
		return NewInvalidParamsResponse(req.ID, err.Error())
	}

	return NewEmptySuccessResponse(req.ID)
}

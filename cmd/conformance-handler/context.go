package main

import (
	"encoding/json"
	"fmt"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// handleContextCreate creates a context with specified chain parameters
func handleContextCreate(registry *Registry, req Request) (Response, error) {
	var params struct {
		ChainParameters struct {
			ChainType string `json:"chain_type"`
		} `json:"chain_parameters"`
		Notifications       *RefObject `json:"notifications,omitempty"`
		ValidationInterface *RefObject `json:"validation_interface,omitempty"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

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
		return Response{}, fmt.Errorf("unknown chain_type: %s", params.ChainParameters.ChainType)
	}

	opts := []kernel.ContextOption{kernel.WithChainType(chainType)}

	if params.Notifications != nil {
		iface, err := registry.GetNotificationCallbacksInterface(params.Notifications.Ref)
		if err != nil {
			return Response{}, err
		}
		opts = append(opts, kernel.WithNotifications(iface))
	}

	if params.ValidationInterface != nil {
		iface, err := registry.GetValidationCallbacksInterface(params.ValidationInterface.Ref)
		if err != nil {
			return Response{}, err
		}
		opts = append(opts, kernel.WithValidationInterface(iface))
	}

	ctx, err := kernel.NewContext(opts...)
	if err != nil {
		return NewEmptyErrorResponse(), nil
	}

	registry.Store(req.Ref, ctx)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleContextDestroy destroys a context
func handleContextDestroy(registry *Registry, req Request) (Response, error) {
	var params struct {
		Context RefObject `json:"context"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if err := registry.Destroy(params.Context.Ref); err != nil {
		return Response{}, err
	}

	return NewEmptySuccessResponse(), nil
}

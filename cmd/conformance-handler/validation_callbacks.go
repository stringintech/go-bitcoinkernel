package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// ValidationCallbacksInterface implements kernel.ValidationInterfaceCallbacks and queues invocation records.
type ValidationCallbacksInterface struct {
	callbackQueue
	enabledCallbacks map[string]bool
}

func (iface *ValidationCallbacksInterface) enabled(name string) bool {
	return iface.enabledCallbacks[name]
}

func validationModeString(m kernel.ValidationMode) string {
	switch m {
	case kernel.ValidationStateValid:
		return "btck_ValidationMode_VALID"
	case kernel.ValidationStateInvalid:
		return "btck_ValidationMode_INVALID"
	case kernel.ValidationStateError:
		return "btck_ValidationMode_INTERNAL_ERROR"
	default:
		return fmt.Sprintf("unknown(%d)", m)
	}
}

func (iface *ValidationCallbacksInterface) BlockChecked(block *kernel.Block, state *kernel.BlockValidationStateView) {
	if !iface.enabled("btck_ValidationInterfaceBlockChecked") {
		return
	}
	// state is a view into a stack-local; must copy before the triggering call returns
	iface.indexedRecord(func(n int) map[string]any {
		return map[string]any{
			"callback": "btck_ValidationInterfaceBlockChecked",
			"block":    iface.ref(fmt.Sprintf("$%s_%d_btck_ValidationInterfaceBlockChecked_block", iface.ifaceRef, n), block),
			"state":    iface.ref(fmt.Sprintf("$%s_%d_btck_ValidationInterfaceBlockChecked_state", iface.ifaceRef, n), state.Copy()),
		}
	})
}

func (iface *ValidationCallbacksInterface) PoWValidBlock(block *kernel.Block, entry *kernel.BlockTreeEntry) {
	if !iface.enabled("btck_ValidationInterfacePoWValidBlock") {
		return
	}
	iface.indexedRecord(func(n int) map[string]any {
		return map[string]any{
			"callback": "btck_ValidationInterfacePoWValidBlock",
			"block":    iface.ref(fmt.Sprintf("$%s_%d_btck_ValidationInterfacePoWValidBlock_block", iface.ifaceRef, n), block),
			"entry":    iface.ref(fmt.Sprintf("$%s_%d_btck_ValidationInterfacePoWValidBlock_entry", iface.ifaceRef, n), entry),
		}
	})
}

func (iface *ValidationCallbacksInterface) BlockConnected(block *kernel.Block, entry *kernel.BlockTreeEntry) {
	if !iface.enabled("btck_ValidationInterfaceBlockConnected") {
		return
	}
	iface.indexedRecord(func(n int) map[string]any {
		return map[string]any{
			"callback": "btck_ValidationInterfaceBlockConnected",
			"block":    iface.ref(fmt.Sprintf("$%s_%d_btck_ValidationInterfaceBlockConnected_block", iface.ifaceRef, n), block),
			"entry":    iface.ref(fmt.Sprintf("$%s_%d_btck_ValidationInterfaceBlockConnected_entry", iface.ifaceRef, n), entry),
		}
	})
}

func (iface *ValidationCallbacksInterface) BlockDisconnected(block *kernel.Block, entry *kernel.BlockTreeEntry) {
	if !iface.enabled("btck_ValidationInterfaceBlockDisconnected") {
		return
	}
	iface.indexedRecord(func(n int) map[string]any {
		return map[string]any{
			"callback": "btck_ValidationInterfaceBlockDisconnected",
			"block":    iface.ref(fmt.Sprintf("$%s_%d_btck_ValidationInterfaceBlockDisconnected_block", iface.ifaceRef, n), block),
			"entry":    iface.ref(fmt.Sprintf("$%s_%d_btck_ValidationInterfaceBlockDisconnected_entry", iface.ifaceRef, n), entry),
		}
	})
}

func handleBlockValidationStateGetValidationMode(registry *Registry, req Request) (Response, error) {
	var params struct {
		State RefObject `json:"state"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	s, err := registry.GetBlockValidationState(params.State.Ref)
	if err != nil {
		return Response{}, err
	}

	return NewSuccessResponse(validationModeString(s.ValidationMode())), nil
}

func handleBlockValidationStateDestroy(registry *Registry, req Request) (Response, error) {
	var params struct {
		State RefObject `json:"state"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if err := registry.Destroy(params.State.Ref); err != nil {
		return Response{}, err
	}

	return NewEmptySuccessResponse(), nil
}

func handleValidationInterfaceCallbacksCreate(registry *Registry, req Request) (Response, error) {
	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	var params struct {
		Callbacks []string `json:"callbacks"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}
	if len(params.Callbacks) == 0 {
		return Response{}, fmt.Errorf("callbacks is required and must not be empty")
	}

	enabledCallbacks := make(map[string]bool, len(params.Callbacks))
	for _, name := range params.Callbacks {
		enabledCallbacks[name] = true
	}

	ifaceRef := strings.TrimPrefix(req.Ref, "$")
	iface := &ValidationCallbacksInterface{
		callbackQueue:    callbackQueue{ifaceRef: ifaceRef, registry: registry},
		enabledCallbacks: enabledCallbacks,
	}
	registry.Store(req.Ref, iface)
	return NewSuccessResponseWithRef(req.Ref), nil
}

func handleValidationCallbacksDrain(registry *Registry, req Request) (Response, error) {
	var params struct {
		Interface RefObject `json:"interface"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	iface, err := registry.GetValidationCallbacksInterface(params.Interface.Ref)
	if err != nil {
		return Response{}, err
	}

	return NewSuccessResponse(iface.drain()), nil
}

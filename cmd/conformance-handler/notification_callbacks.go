package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// NotificationCallbacksInterface implements kernel.NotificationCallbacks and queues invocation records.
type NotificationCallbacksInterface struct {
	callbackQueue
	enabledCallbacks map[string]bool
}

func (iface *NotificationCallbacksInterface) enabled(name string) bool {
	return iface.enabledCallbacks[name]
}

func warningString(w kernel.Warning) string {
	switch w {
	case kernel.WarningUnknownNewRulesActivated:
		return "btck_Warning_UNKNOWN_NEW_RULES_ACTIVATED"
	case kernel.WarningLargeWorkInvalidChain:
		return "btck_Warning_LARGE_WORK_INVALID_CHAIN"
	default:
		return fmt.Sprintf("unknown(%d)", w)
	}
}

func syncStateString(s kernel.SynchronizationState) string {
	switch s {
	case kernel.SyncStatePostInit:
		return "btck_SynchronizationState_POST_INIT"
	case kernel.SyncStateInitDownload:
		return "btck_SynchronizationState_INIT_DOWNLOAD"
	case kernel.SyncStateInitReindex:
		return "btck_SynchronizationState_INIT_REINDEX"
	default:
		return fmt.Sprintf("unknown(%d)", s)
	}
}

func (iface *NotificationCallbacksInterface) BlockTip(state kernel.SynchronizationState, entry *kernel.BlockTreeEntry, progress float64) {
	if !iface.enabled("btck_NotifyBlockTip") {
		return
	}
	iface.indexedRecord(func(n int) map[string]any {
		return map[string]any{
			"callback":              "btck_NotifyBlockTip",
			"state":                 syncStateString(state),
			"entry":                 iface.ref(fmt.Sprintf("$%s_%d_btck_NotifyBlockTip_entry", iface.ifaceRef, n), entry),
			"verification_progress": progress,
		}
	})
}

func (iface *NotificationCallbacksInterface) HeaderTip(state kernel.SynchronizationState, height int64, timestamp int64, presync bool) {
	if !iface.enabled("btck_NotifyHeaderTip") {
		return
	}
	iface.record(map[string]any{
		"callback":  "btck_NotifyHeaderTip",
		"state":     syncStateString(state),
		"height":    height,
		"timestamp": timestamp,
		"presync":   presync,
	})
}

func (iface *NotificationCallbacksInterface) Progress(title string, percent int, resumable bool) {
	if !iface.enabled("btck_NotifyProgress") {
		return
	}
	iface.record(map[string]any{
		"callback":  "btck_NotifyProgress",
		"title":     title,
		"percent":   percent,
		"resumable": resumable,
	})
}

func (iface *NotificationCallbacksInterface) WarningSet(warning kernel.Warning, message string) {
	if !iface.enabled("btck_NotifyWarningSet") {
		return
	}
	iface.record(map[string]any{
		"callback": "btck_NotifyWarningSet",
		"warning":  warningString(warning),
		"message":  message,
	})
}

func (iface *NotificationCallbacksInterface) WarningUnset(warning kernel.Warning) {
	if !iface.enabled("btck_NotifyWarningUnset") {
		return
	}
	iface.record(map[string]any{
		"callback": "btck_NotifyWarningUnset",
		"warning":  warningString(warning),
	})
}

func (iface *NotificationCallbacksInterface) FlushError(message string) {
	if !iface.enabled("btck_NotifyFlushError") {
		return
	}
	iface.record(map[string]any{
		"callback": "btck_NotifyFlushError",
		"message":  message,
	})
}

func (iface *NotificationCallbacksInterface) FatalError(message string) {
	if !iface.enabled("btck_NotifyFatalError") {
		return
	}
	iface.record(map[string]any{
		"callback": "btck_NotifyFatalError",
		"message":  message,
	})
}

func handleNotificationCallbacksCreate(registry *Registry, req Request) (Response, error) {
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
	iface := &NotificationCallbacksInterface{
		callbackQueue:    callbackQueue{ifaceRef: ifaceRef, registry: registry},
		enabledCallbacks: enabledCallbacks,
	}
	registry.Store(req.Ref, iface)
	return NewSuccessResponseWithRef(req.Ref), nil
}

func handleNotificationCallbacksDrain(registry *Registry, req Request) (Response, error) {
	var params struct {
		Interface RefObject `json:"interface"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	iface, err := registry.GetNotificationCallbacksInterface(params.Interface.Ref)
	if err != nil {
		return Response{}, err
	}

	return NewSuccessResponse(iface.drain()), nil
}

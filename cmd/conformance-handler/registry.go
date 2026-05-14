package main

import (
	"fmt"
	"os"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// Registry stores named references to objects created during the test session.
// Objects remain alive throughout the handler's lifetime unless explicitly destroyed.
type Registry struct {
	objects map[string]interface{}
	order   []string // Tracks insertion order for proper cleanup (newest to oldest)
}

// NewRegistry creates a new empty registry
func NewRegistry() *Registry {
	return &Registry{
		objects: make(map[string]interface{}),
		order:   make([]string, 0),
	}
}

// Store stores an object under the given reference name
func (r *Registry) Store(ref string, obj interface{}) {
	// Check if object already exists
	if _, ok := r.objects[ref]; ok {
		// Cleanup the old object before replacing
		_ = r.Destroy(ref)
	}
	r.order = append(r.order, ref)
	r.objects[ref] = obj
}

// GetContext retrieves a context by reference name
func (r *Registry) GetContext(ref string) (*kernel.Context, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	ctx, ok := obj.(*kernel.Context)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a Context (got %T)", ref, obj)
	}
	return ctx, nil
}

// GetChainstateManager retrieves a chainstate manager by reference name
func (r *Registry) GetChainstateManager(ref string) (*ChainstateManagerState, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	csm, ok := obj.(*ChainstateManagerState)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a ChainstateManager (got %T)", ref, obj)
	}
	return csm, nil
}

// GetChain retrieves a chain by reference name
func (r *Registry) GetChain(ref string) (*kernel.Chain, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	chain, ok := obj.(*kernel.Chain)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a Chain (got %T)", ref, obj)
	}
	return chain, nil
}

// GetBlock retrieves a block by reference name
func (r *Registry) GetBlock(ref string) (*kernel.Block, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	block, ok := obj.(*kernel.Block)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a Block (got %T)", ref, obj)
	}
	return block, nil
}

// GetScriptPubkey retrieves a script pubkey by reference name
func (r *Registry) GetScriptPubkey(ref string) (kernel.ScriptPubkeyLike, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	spk, ok := obj.(kernel.ScriptPubkeyLike)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a ScriptPubkey (got %T)", ref, obj)
	}
	return spk, nil
}

// GetTransaction retrieves a transaction by reference name
func (r *Registry) GetTransaction(ref string) (kernel.TransactionLike, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	tx, ok := obj.(kernel.TransactionLike)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a Transaction (got %T)", ref, obj)
	}
	return tx, nil
}

// GetTransactionOutput retrieves a transaction output by reference name
func (r *Registry) GetTransactionOutput(ref string) (kernel.TransactionOutputLike, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	txOut, ok := obj.(kernel.TransactionOutputLike)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a TransactionOutput (got %T)", ref, obj)
	}
	return txOut, nil
}

// GetPrecomputedTransactionData retrieves precomputed transaction data by reference name
func (r *Registry) GetPrecomputedTransactionData(ref string) (*kernel.PrecomputedTransactionData, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	ptd, ok := obj.(*kernel.PrecomputedTransactionData)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a PrecomputedTransactionData (got %T)", ref, obj)
	}
	return ptd, nil
}

// GetTransactionOutPoint retrieves a transaction out point by reference name.
func (r *Registry) GetTransactionOutPoint(ref string) (kernel.TransactionOutPointLike, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	op, ok := obj.(kernel.TransactionOutPointLike)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a TransactionOutPoint (got %T)", ref, obj)
	}
	return op, nil
}

// GetTxid retrieves a txid by reference name.
func (r *Registry) GetTxid(ref string) (kernel.TxidLike, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	txid, ok := obj.(kernel.TxidLike)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a Txid (got %T)", ref, obj)
	}
	return txid, nil
}

// GetTransactionInput retrieves a transaction input by reference name.
func (r *Registry) GetTransactionInput(ref string) (kernel.TransactionInputLike, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	ti, ok := obj.(kernel.TransactionInputLike)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a TransactionInput (got %T)", ref, obj)
	}
	return ti, nil
}

// GetChainParameters retrieves chain parameters by reference name
func (r *Registry) GetChainParameters(ref string) (*kernel.ChainParameters, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	cp, ok := obj.(*kernel.ChainParameters)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a ChainParameters (got %T)", ref, obj)
	}
	return cp, nil
}

// GetBlockHash retrieves a block hash by reference name.
func (r *Registry) GetBlockHash(ref string) (kernel.BlockHashLike, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	bh, ok := obj.(kernel.BlockHashLike)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a BlockHash (got %T)", ref, obj)
	}
	return bh, nil
}

// GetBlockHeader retrieves a block header by reference name
func (r *Registry) GetBlockHeader(ref string) (*kernel.BlockHeader, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	bh, ok := obj.(*kernel.BlockHeader)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a BlockHeader (got %T)", ref, obj)
	}
	return bh, nil
}

// GetNotificationCallbacksInterface retrieves a notification callbacks interface by reference name
func (r *Registry) GetNotificationCallbacksInterface(ref string) (*NotificationCallbacksInterface, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	iface, ok := obj.(*NotificationCallbacksInterface)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a NotificationCallbacksInterface (got %T)", ref, obj)
	}
	return iface, nil
}

// GetBlockValidationState retrieves a block validation state by reference name
func (r *Registry) GetBlockValidationState(ref string) (*kernel.BlockValidationState, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	s, ok := obj.(*kernel.BlockValidationState)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a BlockValidationState (got %T)", ref, obj)
	}
	return s, nil
}

// GetValidationCallbacksInterface retrieves a validation callbacks interface by reference name
func (r *Registry) GetValidationCallbacksInterface(ref string) (*ValidationCallbacksInterface, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	iface, ok := obj.(*ValidationCallbacksInterface)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a ValidationCallbacksInterface (got %T)", ref, obj)
	}
	return iface, nil
}

// GetBlockTreeEntry retrieves a block tree entry by reference name
func (r *Registry) GetBlockTreeEntry(ref string) (*kernel.BlockTreeEntry, error) {
	obj, ok := r.objects[ref]
	if !ok {
		return nil, fmt.Errorf("reference not found: %s", ref)
	}
	entry, ok := obj.(*kernel.BlockTreeEntry)
	if !ok {
		return nil, fmt.Errorf("reference %s is not a BlockTreeEntry (got %T)", ref, obj)
	}
	return entry, nil
}

// Destroy removes and destroys a single object from the registry by reference name
func (r *Registry) Destroy(ref string) error {
	obj, ok := r.objects[ref]
	if !ok {
		return fmt.Errorf("reference not found: %s", ref)
	}

	// Destroy the object
	r.destroyObject(obj)

	// Remove from registry
	delete(r.objects, ref)

	// Remove from order slice
	for i, name := range r.order {
		if name == ref {
			r.order = append(r.order[:i], r.order[i+1:]...)
			break
		}
	}

	return nil
}

// Cleanup destroys all objects in the registry and clears all references
// Objects are destroyed in reverse order (newest to oldest) to handle dependencies
func (r *Registry) Cleanup() {
	// Destroy objects in reverse order (newest to oldest)
	for i := len(r.order) - 1; i >= 0; i-- {
		ref := r.order[i]
		if obj, ok := r.objects[ref]; ok {
			r.destroyObject(obj)
		}
	}

	// Clear everything
	r.objects = make(map[string]interface{})
	r.order = nil
}

// destroyObject releases owned resources. View-like objects without Destroy are ignored.
func (r *Registry) destroyObject(obj interface{}) {
	if v, ok := obj.(interface{ Destroy() }); ok {
		v.Destroy()
	}
}

// ChainstateManagerState holds the chainstate manager and its dependencies
type ChainstateManagerState struct {
	Manager *kernel.ChainstateManager
	TempDir string
}

// Destroy releases all resources held by the chainstate manager state
func (c *ChainstateManagerState) Destroy() {
	if c.Manager != nil {
		c.Manager.Destroy()
		c.Manager = nil
	}

	// Remove temp directory if it exists
	if c.TempDir != "" {
		_ = os.RemoveAll(c.TempDir)
		c.TempDir = ""
	}
}

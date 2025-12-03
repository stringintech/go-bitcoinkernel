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

// destroyObject performs cleanup on a single object based on its type
func (r *Registry) destroyObject(obj interface{}) {
	switch v := obj.(type) {
	case *kernel.Context:
		if v != nil {
			v.Destroy()
		}
	case *ChainstateManagerState:
		if v != nil {
			v.Cleanup()
		}
	case *kernel.Block:
		if v != nil {
			v.Destroy()
		}
		// Chain and BlockTreeEntry don't need explicit cleanup
	}
}

// ChainstateManagerState holds the chainstate manager and its dependencies
type ChainstateManagerState struct {
	Manager *kernel.ChainstateManager
	TempDir string
}

// Cleanup releases all resources held by the chainstate manager state
func (c *ChainstateManagerState) Cleanup() {
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

// Package kernel provides Go bindings for the Bitcoin Core kernel library.
//
// # Memory Management
//
// Objects that provide a Destroy() method hold owned references to underlying C resources.
// These resources should be freed immediately by calling Destroy() when no longer needed:
//
//	block, err := NewBlock(rawBlockBytes)
//	if err != nil {
//	    return err
//	}
//	defer block.Destroy()
//
// If Destroy() is not called explicitly, the garbage collector will eventually free
// the resources automatically via finalizers. However, relying on finalizers may delay
// resource cleanup and is not recommended for long-running programs or when working
// with many objects.
//
// # Owned and View Types
//
// Some kernel objects have both an owned type and a view type, such as Transaction
// and TransactionView. Owned types hold a reference to an underlying C resource and
// provide Destroy. View types are non-owned pointers returned from another object or
// callback, and remain valid only while the object that produced them remains valid.
//
// Methods shared by an owned type and its view type are exposed through sealed
// interfaces named with the Like suffix, such as TransactionLike and TxidLike. These
// interfaces are implemented by the package's owned and view types, but cannot be
// implemented by external packages because they include an unexported pointer getter.
// This keeps APIs type-safe around kernel C pointers while still allowing callers to
// accept either owned or view values.
package kernel

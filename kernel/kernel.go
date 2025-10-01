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
package kernel

/*
#cgo CFLAGS: -I../depend/bitcoin/src
#cgo LDFLAGS: -L../depend/bitcoin/build/lib -lbitcoinkernel -Wl,-rpath,${SRCDIR}/../depend/bitcoin/build/lib
*/
import "C"

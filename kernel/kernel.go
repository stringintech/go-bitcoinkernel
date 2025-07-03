// Package kernel provides Go bindings for the Bitcoin Core kernel library.
//
// Resource Management:
// All types that implement the cManagedResource interface follow a dual cleanup pattern:
//  1. Explicit cleanup via Destroy() methods (preferred)
//  2. Finalizers as safety net for forgotten cleanup
package kernel

/*
#cgo CFLAGS: -I../depend/bitcoin/src
#cgo LDFLAGS: -L../depend/bitcoin/build/lib -lbitcoinkernel -Wl,-rpath,${SRCDIR}/../depend/bitcoin/build/lib
*/
import "C"

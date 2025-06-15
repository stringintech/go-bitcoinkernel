// Package kernel provides Go bindings for the Bitcoin Core kernel library.
// This package offers an interface to Bitcoin Core's kernel functionality,
//
// Resource Management:
// All types that wrap C resources (Context, ContextOptions, ChainParameters, etc.)
// follow a dual cleanup pattern:
//  1. Explicit cleanup via Close() methods (preferred)
//  2. Finalizers as safety net for forgotten cleanup
//
// Always call Close() explicitly, preferably with defer:
//
//	ctx, err := NewDefaultContext()
//	if err != nil { return err }
//	defer ctx.Close()
package kernel

/*
#cgo CFLAGS: -I../depend/bitcoin/src
#cgo LDFLAGS: -L../depend/bitcoin/build/lib -lbitcoinkernel -Wl,-rpath,${SRCDIR}/../depend/bitcoin/build/lib
*/
import "C"

// ReverseBytes reverses bytes for display (Bitcoin hashes are displayed in reverse order)
func ReverseBytes(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[len(data)-1-i] = b
	}
	return result
}

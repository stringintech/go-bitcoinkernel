//go:build unix

package kernel

/*
#cgo CFLAGS: -I../depend/bitcoin/src
#cgo LDFLAGS: -L../depend/bitcoin/build/lib -lbitcoinkernel -Wl,-rpath,${SRCDIR}/../depend/bitcoin/build/lib
*/
import "C"

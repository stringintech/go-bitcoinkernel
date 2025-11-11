//go:build unix

package kernel

/*
#cgo CFLAGS: -I../depend/bitcoin/install/include -DBITCOINKERNEL_STATIC
#cgo darwin LDFLAGS: -L../depend/bitcoin/install/lib -lbitcoinkernel -lc++
#cgo !darwin LDFLAGS: -L../depend/bitcoin/install/lib -lbitcoinkernel -lstdc++ -lm
*/
import "C"

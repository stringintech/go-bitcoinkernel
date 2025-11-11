//go:build windows

package kernel

/*
#cgo CFLAGS: -I../depend/bitcoin/install/include -DBITCOINKERNEL_STATIC
#cgo LDFLAGS: -L../depend/bitcoin/install/lib -lbitcoinkernel -lstdc++ -lbcrypt -lshell32
*/
import "C"

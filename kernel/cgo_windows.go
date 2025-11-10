//go:build windows

package kernel

/*
#cgo CFLAGS: -I../depend/bitcoin/src
#cgo LDFLAGS: -L../depend/bitcoin/build/bin/RelWithDebInfo -lbitcoinkernel -lbcrypt -lshell32
*/
import "C"

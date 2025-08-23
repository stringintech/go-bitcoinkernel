package bench

import (
	"encoding/hex"
	"github.com/stringintech/go-bitcoinkernel/kernel"
	"os"
	"strings"
	"testing"
)

func loadTestBlock(tb testing.TB) *kernel.Block {
	tb.Helper()

	// Load a recent mainnet block
	hexBytes, err := os.ReadFile("../data/mainnet/block911451.txt")
	if err != nil {
		tb.Fatal(err)
	}
	hexString := strings.TrimSpace(string(hexBytes))

	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		tb.Fatal(err)
	}

	block, err := kernel.NewBlockFromRaw(bytes)
	if err != nil {
		tb.Fatal(err)
	}

	return block
}

func BenchmarkComparison(b *testing.B) {
	block := loadTestBlock(b)
	defer block.Destroy()

	b.Run("Bytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := block.Bytes(); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("PreAllocBytes", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := block.PreAllocBytes(); err != nil {
				b.Fatal(err)
			}
		}
	})
}

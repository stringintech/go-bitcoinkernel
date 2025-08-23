# Benchmarks

## Block Serialization Benchmark

**File:** `block_bytes_bench_test.go`

**What it compares:**
Two different approaches to serializing blocks to bytes:

1. **Bytes()** - Uses growing slice with dynamic allocations
2. **PreAllocBytes()** - Pre-allocates slice using `btck_block_get_serialize_size` C API

### Running the Benchmark

```bash
go test -bench=BenchmarkComparison -benchmem -benchtime=10x ./bench/
```

### Sample Results

```
BenchmarkComparison/Bytes         	      10	   7074350 ns/op	10590606 B/op	      43 allocs/op
BenchmarkComparison/PreAllocBytes 	      10	   6973412 ns/op	 1802627 B/op	       9 allocs/op
```

- **Memory efficiency**: ~83% reduction in total allocations (10.6MB → 1.8MB)
- **Allocation count**: ~79% reduction (43 → 9 allocations) 
- **CPU performance**: ~1.4% improvement (7.07ms → 6.97ms)
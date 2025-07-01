# Bitcoin Kernel Go Wrapper - TODO

This document lists the remaining C API functions and data structures from [`bitcoinkernel.h`](./depend/bitcoin/src/kernel/bitcoinkernel.h) that haven't been wrapped yet in the Go kernel package.

## Missing Data Structures

- **`kernel_BlockPointer`** - Non-owned block pointers (from callbacks)

## Missing Functions by Category

### Script Operations
- [ ] `kernel_verify_script()` - **Script verification (IMPORTANT!)**

### Block Operations (Additional)
- [ ] `kernel_block_pointer_get_hash()` - Get hash from block pointer
- [ ] `kernel_copy_block_pointer_data()` - Copy data from block pointer

### Callback Support
- [ ] **Notification callbacks** - Full integration of kernel notification system
- [ ] **Validation interface callbacks** - Block validation event handling

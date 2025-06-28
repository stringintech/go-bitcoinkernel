# Bitcoin Kernel Go Wrapper - TODO

This document lists the remaining C API functions and data structures from [`bitcoinkernel.h`](./depend/bitcoin/src/kernel/bitcoinkernel.h) that haven't been wrapped yet in the Go kernel package.

## Missing Data Structures

### Core Transaction Types
- **`kernel_BlockPointer`** - Non-owned block pointers (from callbacks)
- **`kernel_BlockUndo`** - Block undo data operations

## Missing Functions by Category

### Script Operations
- [ ] `kernel_verify_script()` - **Script verification (IMPORTANT!)**

### Block Operations (Additional)
- [ ] `kernel_block_pointer_get_hash()` - Get hash from block pointer
- [ ] `kernel_copy_block_pointer_data()` - Copy data from block pointer  

### Block Undo Operations
- [ ] `kernel_read_block_undo_from_disk()` - Read undo data from disk
- [ ] `kernel_block_undo_size()` - Get number of transactions in undo data
- [ ] `kernel_get_transaction_undo_size()` - Get output count per transaction
- [ ] `kernel_get_undo_output_height_by_index()` - Get output block height
- [ ] `kernel_get_undo_output_by_index()` - Get specific undo output
- [ ] `kernel_block_undo_destroy()` - Cleanup undo data

### Callback Support
- [ ] **Notification callbacks** - Full integration of kernel notification system
- [ ] **Validation interface callbacks** - Block validation event handling

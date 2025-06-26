# Bitcoin Kernel Go Wrapper - TODO

This document lists the remaining C API functions and data structures from `bitcoinkernel.h` that haven't been wrapped yet in the Go kernel package.

## Current Status

✅ **Implemented:**
- Chain parameters and context management
- Block and block index operations
- Chainstate manager functionality
- Basic logging support

❌ **Missing:**
- Transaction processing pipeline
- Script validation system
- Block undo data handling
- Complete callback system integration

## Missing Data Structures

### Core Transaction Types
- **`kernel_TransactionOutput`** - Transaction output operations
- **`kernel_BlockPointer`** - Non-owned block pointers (from callbacks)
- **`kernel_BlockUndo`** - Block undo data operations

## Missing Functions by Category

### Script Operations
- [ ] `kernel_verify_script()` - **Script verification (IMPORTANT!)**

### Transaction Output Operations
- [ ] `kernel_transaction_output_create()` - Create transaction output
- [ ] `kernel_transaction_output_destroy()` - Cleanup transaction output
- [ ] `kernel_copy_script_pubkey_from_output()` - Extract script from output
- [ ] `kernel_get_transaction_output_amount()` - Get output amount

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

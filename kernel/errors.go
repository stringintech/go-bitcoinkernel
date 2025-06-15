package kernel

import "errors"

var (
	// ErrChainParametersCreation indicates failure to create chain parameters
	ErrChainParametersCreation = errors.New("failed to create chain parameters")
	// ErrContextOptionsCreation indicates failure to create context options
	ErrContextOptionsCreation = errors.New("failed to create context options")
	// ErrContextCreation indicates failure to create kernel context
	ErrContextCreation = errors.New("failed to create kernel context")
	// ErrInvalidChainType indicates an invalid chain type was provided
	ErrInvalidChainType = errors.New("invalid chain type")
	// ErrChainstateManagerOptionsCreation indicates failure to create chainstate manager options
	ErrChainstateManagerOptionsCreation = errors.New("failed to create chainstate manager options")
	// ErrChainstateManagerCreation indicates failure to create chainstate manager
	ErrChainstateManagerCreation = errors.New("failed to create chainstate manager")
	// ErrBlockCreation indicates failure to create block from raw data
	ErrBlockCreation = errors.New("failed to create block from raw data")
	// ErrInvalidBlockData indicates invalid block data was provided
	ErrInvalidBlockData = errors.New("invalid block data")
	// ErrInvalidBlock indicates block is invalid or nil
	ErrInvalidBlock = errors.New("invalid block")
	// ErrInvalidBlockIndex indicates block index is invalid or nil
	ErrInvalidBlockIndex = errors.New("invalid block index")
	// ErrHashCalculation indicates failure to calculate hash
	ErrHashCalculation = errors.New("failed to calculate hash")
	// ErrBlockDataCopy indicates failure to copy block data
	ErrBlockDataCopy = errors.New("failed to copy block data")
	// ErrBlockProcessing indicates failure to process block
	ErrBlockProcessing = errors.New("failed to process block")
	// ErrBlockRead indicates failure to read block from disk
	ErrBlockRead = errors.New("failed to read block from disk")
)

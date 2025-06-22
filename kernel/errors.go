package kernel

import "errors"

var (
	ErrChainParametersCreation = errors.New("failed to create chain parameters")
	ErrContextOptionsCreation  = errors.New("failed to create context options")
	ErrContextCreation         = errors.New("failed to create kernel context")
	ErrInvalidChainType        = errors.New("invalid chain type")
	ErrBlockCreation           = errors.New("failed to create block from raw data")
	ErrInvalidBlockData        = errors.New("invalid block data")
	ErrInvalidBlock            = errors.New("invalid block")
	ErrHashCalculation         = errors.New("failed to calculate hash")
	ErrBlockDataCopy           = errors.New("failed to copy block data")
)

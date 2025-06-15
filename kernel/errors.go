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
)

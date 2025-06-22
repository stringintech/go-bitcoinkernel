package kernel

import "errors"

var (
	ErrChainParametersCreation = errors.New("failed to create chain parameters")
	ErrContextOptionsCreation  = errors.New("failed to create context options")
	ErrContextCreation         = errors.New("failed to create kernel context")
	ErrInvalidChainType        = errors.New("invalid chain type")
)

package kernel

import (
	"errors"
	"testing"
)

func TestNewChainParameters(t *testing.T) {
	tests := []struct {
		name      string
		chainType ChainType
		wantErr   bool
		errType   error
	}{
		{
			name:      "Mainnet chain parameters",
			chainType: ChainTypeMainnet,
			wantErr:   false,
		},
		{
			name:      "Testnet chain parameters",
			chainType: ChainTypeTestnet,
			wantErr:   false,
		},
		{
			name:      "Testnet4 chain parameters",
			chainType: ChainTypeTestnet4,
			wantErr:   false,
		},
		{
			name:      "Signet chain parameters",
			chainType: ChainTypeSignet,
			wantErr:   false,
		},
		{
			name:      "Regtest chain parameters",
			chainType: ChainTypeRegtest,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp, err := NewChainParameters(tt.chainType)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewChainParameters() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("NewChainParameters() error = %v, want %v", err, tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("NewChainParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if cp == nil {
				t.Error("NewChainParameters() returned nil ChainParameters")
				return
			}

			if cp.ptr == nil {
				t.Error("ChainParameters has nil pointer")
			}

			// Clean up
			cp.Destroy()
		})
	}
}

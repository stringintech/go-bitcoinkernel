package kernel

import (
	"testing"
)

func TestNewContext(t *testing.T) {
	tests := []struct {
		name        string
		setupOption func() *ContextOptions
		wantErr     bool
		errType     error
	}{
		{
			name: "Valid context options",
			setupOption: func() *ContextOptions {
				opts := NewContextOptions()
				params, err := NewChainParameters(ChainTypeMainnet)
				if err != nil {
					t.Fatalf("Failed to create chain parameters: %v", err)
				}
				defer params.Destroy()
				opts.SetChainParams(params)
				return opts
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.setupOption()

			ctx, err := NewContext(opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewContext() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errType != nil && err.Error() != tt.errType.Error() {
					t.Errorf("NewContext() error = %v, want %v", err, tt.errType)
				}
				return
			}

			if err != nil {
				t.Errorf("NewContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if ctx == nil {
				t.Error("NewContext() returned nil Context")
				return
			}

			if ctx.ptr == nil {
				t.Error("Context has nil pointer")
			}

			// Clean up
			ctx.Destroy()
		})
	}
}

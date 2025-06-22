package kernel

import (
	"errors"
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
				opts, err := NewContextOptions()
				if err != nil {
					t.Fatalf("Failed to create context options: %v", err)
				}
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
		{
			name: "Nil context options",
			setupOption: func() *ContextOptions {
				return nil
			},
			wantErr: true,
			errType: errors.New("context options cannot be nil"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.setupOption()
			if opts != nil {
				defer opts.Destroy()
			}

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

func TestNewDefaultContext(t *testing.T) {
	ctx, err := NewDefaultContext()

	if err != nil {
		t.Errorf("NewDefaultContext() error = %v, want nil", err)
		return
	}

	if ctx == nil {
		t.Error("NewDefaultContext() returned nil Context")
		return
	}

	if ctx.ptr == nil {
		t.Error("Context has nil pointer")
	}
}

func TestContextInterrupt(t *testing.T) {
	ctx, err := NewDefaultContext()
	if err != nil {
		t.Fatalf("Failed to create default context: %v", err)
	}
	defer ctx.Destroy()

	// Test interrupt on valid context
	result := ctx.Interrupt()
	// Result can be true or false, both could be valid
	t.Logf("Context interrupt result: %v", result)

	// Destroy and test interrupt on destroyed context
	ctx.Destroy()
	result = ctx.Interrupt()
	if result {
		t.Error("Interrupt() should return false after context is destroyed")
	}
}

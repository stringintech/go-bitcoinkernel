package kernel

import "testing"

func TestDefaultContext(t *testing.T) {
	ctx, err := NewDefaultContext()
	if err != nil {
		t.Fatalf("NewDefaultContext() error = %v", err)
	}
	defer ctx.Close()

	if !ctx.IsValid() {
		t.Error("Context should be valid after creation")
	}
}

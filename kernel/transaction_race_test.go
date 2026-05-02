//go:build race

// This race-only test shows why owned API structs must not cache their own C
// pointer. When transactionApi.Copy calls the pointer getter, it reads
// handle.ptr, so the race detector can catch Copy running concurrently with
// Destroy, which writes to handle.ptr. If transactionApi cached the C pointer
// instead, Copy would read that cached field and the race detector would not see
// the unsafe shared access to handle.ptr.
package kernel

import (
	"bytes"
	"encoding/hex"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

func TestTransactionHandlePtrRaceDetector(t *testing.T) {
	if os.Getenv("KERNEL_RACE_HELPER") == "1" {
		runTransactionHandlePtrRaceHelper(t)
		return
	}

	// Run the helper in a subprocess because it intentionally triggers a race.
	// GORACE makes the helper exit successfully after printing the first report,
	// so this outer test can assert that the expected report was produced.
	cmd := exec.Command(os.Args[0], "-test.run=TestTransactionHandlePtrRaceDetector", "-test.v")
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "GORACE=") {
			continue
		}
		if strings.HasPrefix(env, "GOMAXPROCS=") {
			continue
		}
		cmd.Env = append(cmd.Env, env)
	}
	cmd.Env = append(cmd.Env,
		"KERNEL_RACE_HELPER=1",
		"GOMAXPROCS=1",
		"GORACE=halt_on_error=1 exitcode=0",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("race helper exited unexpectedly: %v\noutput:\n%s", err, output)
	}
	report := string(output)
	if !bytes.Contains(output, []byte("WARNING: DATA RACE")) {
		t.Fatalf("expected race detector output, got:\n%s", output)
	}
	if !regexp.MustCompile(`(?s)[Rr]ead at .*?(\(\*transactionApi\)\.Copy\(\))`).MatchString(report) {
		t.Fatalf("expected race output to mention transactionApi.Copy, got:\n%s", output)
	}
	if !regexp.MustCompile(`(?s)[Ww]rite at .*?(\(\*handle\)\.Destroy\(\))`).MatchString(report) {
		t.Fatalf("expected race output to mention handle.Destroy, got:\n%s", output)
	}
}

func runTransactionHandlePtrRaceHelper(t *testing.T) {
	rawTransaction, err := hex.DecodeString(coinbaseTxHex)
	if err != nil {
		t.Fatalf("DecodeString() error = %v", err)
	}

	tx, err := NewTransaction(rawTransaction)
	if err != nil {
		t.Fatalf("NewTransaction() error = %v", err)
	}

	// Race repeated public API reads against destroying the same owned value.
	// Copy() is enough to prove the point because it reads handle.ptr before
	// creating another owned Transaction.
	ready := make(chan struct{})
	done := make(chan struct{}, 2)

	go func() {
		// Warm up the copy loop so it is already active before Destroy().
		for i := 0; i < 50; i++ {
			_ = tx.Copy()
			runtime.Gosched()
		}
		close(ready)
		// Keep copying after Destroy() starts to widen the race window.
		for i := 0; i < 1000; i++ {
			func() {
				defer func() {
					_ = recover()
				}()
				_ = tx.Copy()
			}()
			runtime.Gosched()
		}
		done <- struct{}{}
	}()

	go func() {
		<-ready
		// Give the copy loop one more scheduling opportunity, then destroy.
		runtime.Gosched()
		tx.Destroy()
		done <- struct{}{}
	}()

	<-done
	<-done
}

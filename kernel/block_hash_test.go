package kernel

import (
	"testing"
)

func TestBlockHash(t *testing.T) {
	hashBytes := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

	hash := NewBlockHash(hashBytes)
	defer hash.Destroy()

	// Test Bytes()
	newHash := NewBlockHash(hash.Bytes())
	defer newHash.Destroy()

	if hash.Bytes() != newHash.Bytes() {
		t.Errorf("Hash bytes differ: %x != %x", hash.Bytes(), newHash.Bytes())
	}

	// Test Copy()
	copiedHash := hash.Copy()
	defer copiedHash.Destroy()

	if hash.Bytes() != copiedHash.Bytes() {
		t.Errorf("Copied hash bytes differ: %x != %x", hash.Bytes(), copiedHash.Bytes())
	}

	// Test Equals()
	if !hash.Equals(copiedHash) {
		t.Errorf("hash.Equals(copiedHash) = false, want true")
	}

	differentHash := NewBlockHash([32]byte{0xFF})
	defer differentHash.Destroy()

	if hash.Equals(differentHash) {
		t.Errorf("hash.Equals(differentHash) = true, want false")
	}
}

package lib

import (
	"crypto/sha256"
	"testing"
)

func TestXorBytes(t *testing.T) {
	a := []byte{0, 1, 1, 0}
	b := []byte{1, 0, 0, 1}

	result := XorBytes(a, b)

	for k, item := range result {
		if item != 1 {
			t.Fatalf("Invalid byte key %d should be 1 but got %d", k, item)
		}
	}
}

func TestXorUuidsShouldBeSame(t *testing.T) {
	txs := []string{"3f3979f4-7828-435e-ae8e-717f62b8d8b1", "22bed1fb-6d2d-4e2d-ac04-91c6a44012b7", "a8ebbe74-1198-4056-96d8-5e276432fd25"}
	hashes := [][32]byte{sha256.Sum256([]byte(txs[0])), sha256.Sum256([]byte(txs[1])), sha256.Sum256([]byte(txs[2]))}

	result1 := XorBytes(hashes[0][:], hashes[1][:])
	result2 := XorBytes(hashes[1][:], hashes[0][:])

	for key, item := range result1 {
		if item != result2[key] {
			t.Fatalf("Invalid byte key %d should be %d but got %d", key, item, result2[key])
		}
	}
}

func TestXorAssociation(t *testing.T) {
	txs := []string{"3f3979f4-7828-435e-ae8e-717f62b8d8b1", "22bed1fb-6d2d-4e2d-ac04-91c6a44012b7", "a8ebbe74-1198-4056-96d8-5e276432fd25"}
	hashes := [][32]byte{sha256.Sum256([]byte(txs[0])), sha256.Sum256([]byte(txs[1])), sha256.Sum256([]byte(txs[2]))}

	result1 := XorBytes(XorBytes(hashes[0][:], hashes[1][:]), hashes[2][:])
	result2 := XorBytes(XorBytes(hashes[2][:], hashes[1][:]), hashes[0][:])

	for key, item := range result1 {
		if item != result2[key] {
			t.Fatalf("Invalid byte key %d should be %d but got %d", key, item, result2[key])
		}
	}
}

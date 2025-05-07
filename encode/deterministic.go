package encode

import (
	"crypto/sha256"
	"encoding/binary"
	"hash"
)

// DeterministicSecureRandom implements io.Reader to provide deterministic random bytes
type DeterministicSecureRandom struct {
	digest  hash.Hash
	seed    []byte
	counter uint32
}

// NewDeterministicSecureRandom creates a new deterministic secure random number generator
func NewDeterministicSecureRandom(seed []byte) *DeterministicSecureRandom {
	// Create a copy of the seed to prevent modifications
	seedCopy := make([]byte, len(seed))
	copy(seedCopy, seed)

	return &DeterministicSecureRandom{
		digest:  sha256.New(),
		seed:    seedCopy,
		counter: 0,
	}
}

// Read implements io.Reader interface
func (d *DeterministicSecureRandom) Read(p []byte) (n int, err error) {
	// Reset counter to ensure deterministic output
	d.counter = 0

	index := 0
	for index < len(p) {
		d.digest.Reset()
		d.digest.Write(d.seed)

		// Convert counter to bytes
		counterBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(counterBytes, d.counter)
		d.counter++

		d.digest.Write(counterBytes)
		hash := d.digest.Sum(nil)

		// Copy as much as we can from the hash
		toCopy := min(len(hash), len(p)-index)
		copy(p[index:], hash[:toCopy])
		index += toCopy
	}

	return len(p), nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

const (
	dirsnapDir = ".dirsnap"
	logFile    = ".dirsnap/log"
)

func objectsDir() string {
	return filepath.Join(dirsnapDir, "objects")
}

// HashAndStore takes raw bytes, computes their SHA-256 hash, and writes them
// to .dirsnap/objects/<hash>. If an object with that hash already exists, it
// does nothing — this is where deduplication comes from for free: identical
// content always produces the identical hash, so we never store the same
// bytes twice, whether it's two copies of one file or two snapshots that
// share a file.
func HashAndStore(data []byte) (string, error) {
	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])

	path := filepath.Join(objectsDir(), hash)
	if _, err := os.Stat(path); err == nil {
		return hash, nil // already stored, nothing to do
	}

	return hash, os.WriteFile(path, data, 0644)
}

// ReadObject reads back the raw bytes stored under a given hash.
func ReadObject(hash string) ([]byte, error) {
	return os.ReadFile(filepath.Join(objectsDir(), hash))
}

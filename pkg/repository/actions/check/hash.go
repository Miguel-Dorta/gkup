package check

import (
	"crypto/sha256"
	"errors"
	"hash"
	"io"
	"os"
	"strings"
)

func getHash(algorithm string) (hash.Hash, error) {
	switch strings.ToLower(algorithm) {
	case "sha256":
		return sha256.New(), nil
	default:
		return nil, errors.New("hash algorithm unknown")
	}
}

func hashFile(path string, h hash.Hash, buf []byte) ([]byte, error) {
	h.Reset()

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err = io.CopyBuffer(h, f, buf); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}


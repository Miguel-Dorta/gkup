package check

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"golang.org/x/crypto/sha3"
	"hash"
	"io"
	"os"
	"strings"
)

func getHash(algorithm string) (hash.Hash, error) {
	switch strings.ToLower(algorithm) {
	case "sha256":
		return sha256.New(), nil
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha512":
		return sha512.New(), nil
	case "sha3-256":
		return sha3.New256(), nil
	case "sha3-512":
		return sha3.New512(), nil
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


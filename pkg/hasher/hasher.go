package hasher

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg"
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"golang.org/x/crypto/sha3"
	"hash"
	"io"
	"os"
	"strings"
)

// Hasher is a type for making hashing operations
type Hasher struct {
	hash hash.Hash
	buf  []byte
}

// New creates a new Hasher object
func New(algorithm string) (*Hasher, error) {
	var h hash.Hash

	switch strings.ToLower(algorithm) {
	case "sha256":
		h = sha256.New() // Most frequent case the first
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha512":
		h = sha512.New()
	case "sha3-256":
		h = sha3.New256()
	case "sha3-512":
		h = sha3.New512()
	default:
		return nil, fmt.Errorf("hash algorithm %s is not supported", algorithm)
	}

	return &Hasher{
		hash: h,
		buf:  make([]byte, pkg.BufferSize),
	}, nil
}

// HashFile gets and assigns the hash from the files.File provided.
func (h *Hasher) HashFile(f *files.File) error {
	if f.RealPath == "" {
		return fmt.Errorf("undefined RealPath in file with name \"%s\"", f.Name)
	}

	file, err := os.Open(f.RealPath)
	if err != nil {
		return fmt.Errorf("cannot open file \"%s\": %s", f.RealPath, err.Error())
	}
	defer file.Close()

	h.hash.Reset()
	pkg.Log.Debugf("Hashing file %s", f.RealPath)
	if _, err := io.CopyBuffer(h.hash, file, h.buf); err != nil {
		return fmt.Errorf("error hashing file \"%s\": %s", f.RealPath, err.Error())
	}

	f.Hash = h.hash.Sum(nil)
	return nil
}

// HashPath hashes the path provided
func (h *Hasher) HashPath(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file \"%s\": %s", path, err.Error())
	}
	defer f.Close()

	h.hash.Reset()
	pkg.Log.Debugf("Hashing path %s", path)
	if _, err := io.CopyBuffer(h.hash, f, h.buf); err != nil {
		return nil, fmt.Errorf("error hashing file \"%s\": %s", path, err.Error())
	}
	return h.hash.Sum(nil), nil
}

// GetFile returns a files.File object from the path provided
func (h *Hasher) GetFile(path string) (*files.File, error) {
	f, err := files.NewFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot get information of \"%s\": %s", path, err.Error())
	}

	pkg.Log.Debugf("Hashing path %s and returning file", path)
	if f.Hash, err = h.HashPath(path); err != nil {
		return nil, err
	}

	return f, nil
}

// CheckFileIntegrity checks if the file of the path provided matches the info implicit in its name.
// This info follow the files.GetFileFromName specification.
func (h *Hasher) CheckFileIntegrity(path string) error {
	// Get actual info of the path provided
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot get file info of file \"%s\": %s", path, err)
	}

	// Get original info from the name
	f, err := files.GetFileFromName(stat.Name())
	if err != nil {
		return fmt.Errorf("name of file \"%s\" is corrupted", path)
	}

	// Check size
	if f.Size != stat.Size() {
		return fmt.Errorf("sizes don't match in file \"%s\"", path)
	}

	// Hash file and check it
	hash, err := h.HashPath(path)
	if err != nil {
		return fmt.Errorf("cannot get hash of file \"%s\": %s", path, err)
	}
	if !bytes.Equal(f.Hash, hash) {
		return fmt.Errorf("hashes don't match in file \"%s\"", path)
	}

	return nil
}

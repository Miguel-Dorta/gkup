package files

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
)

// File represents an abstraction of a file
type File struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Hash     []byte `json:"hash"`
	RealPath string `json:"-"`
}

// NewFile gets a File object from the path provided without hashing it
func NewFile(path string) (*File, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot get information of \"%s\": %s", path, err.Error())
	}

	return &File{
		Name:     stat.Name(),
		Size:     stat.Size(),
		Hash:     nil,
		RealPath: path,
	}, nil
}

// GetFileFromName gets a File object from a string formed by the file's hash and its size.
// The hash is encoded in hexadecimal notation. The size is written in base 10.
// Both are separated by the character '-'
func GetFileFromName(fileName string) (*File, error) {
	var err error
	f := File{
		Name: fileName,
		Size: -1,
	}

	for i, b := range fileName {
		if b != '-' {
			continue
		}

		if f.Hash, err = hex.DecodeString(fileName[:i]); err != nil {
			return nil, fmt.Errorf("cannot decode hash: %s", err.Error())
		}

		if f.Size, err = strconv.ParseInt(fileName[i+1:], 10, 64); err != nil {
			return nil, fmt.Errorf("cannot parse size: %s", err.Error())
		}
		break
	}

	if f.Hash == nil || f.Size < 0 {
		return nil, errors.New("invalid format")
	}

	return &f, nil
}

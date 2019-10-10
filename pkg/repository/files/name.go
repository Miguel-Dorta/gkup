package files

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func GetDataFromName(name string) (hash []byte, size int64, err error) {
	// Get index of character '-'
	separatorIndex := strings.IndexByte(name, '-')
	if separatorIndex < 0 {
		return nil, -1, errors.New("incorrect file name")
	}

	// Get hash
	hash, err = hex.DecodeString(name[:separatorIndex])
	if err != nil {
		return nil, -1, fmt.Errorf("error decoding hash from name: %w", err)
	}

	// Get size
	size, err = strconv.ParseInt(name[separatorIndex+1:], 10, 64)
	if err != nil {
		return nil, -1, fmt.Errorf("error parsing size from name: %w", err)
	}

	return hash, size, nil
}

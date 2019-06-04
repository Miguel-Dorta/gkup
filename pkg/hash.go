package pkg

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
)

func hashFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file \"%s\": %s", path, err.Error())
	}
	defer f.Close()

	var h hash.Hash
	if HashAlgorithm == "sha1" {
		h = sha1.New()
	} else {
		h = sha256.New()
	}

	// Hash it
	if _, err := io.Copy(h, f); err != nil {
		return nil, fmt.Errorf("error hashing file \"%s\": %s", path, err.Error())
	}

	return h.Sum(nil), nil
}

/*func getIntSliceFromByteSlice(b []byte) []uint64 {
	bLen := len(b)

	// Create int slice
	iSize := bLen / 8
	if (bLen % 8) != 0 {
		iSize++
	}
	intSlice := make([]uint64, iSize)

	// Convert byte slice into int slice
	for i:=0; i<bLen; i+=8 {
		dest := make([]byte, 8)
		end := i + 8
		if bLen < end {
			end = bLen
		}
		copy(dest, b[i:end])
		intSlice[i/8] = binary.BigEndian.Uint64(dest)
	}

	return intSlice
}

func getByteSliceFromIntSlice(intSlice []uint64) []byte {
	result := make([]byte, len(intSlice) * 8)
	for i, x := range intSlice {
		binary.BigEndian.PutUint64(result[i*8 : i*8+8], x)
	}
	return result
}*/

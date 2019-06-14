package pkg

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"golang.org/x/crypto/sha3"
	"golang.org/x/sync/errgroup"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type hasher struct {
	hash hash.Hash
	buf []byte
}

type multiHasher struct {
	workers []*hasher
}

func newMultiHasher(algorithm string, bufferSize, threads int) (*multiHasher, error) {
	var err error

	workers := make([]*hasher, threads)
	for i := range workers{
		workers[i], err = newHasher(algorithm, bufferSize)
		if err != nil {
			return nil, err
		}
	}
	return &multiHasher{workers: workers}, nil
}

func newHasher(algorithm string, bufferSize int) (*hasher, error) {
	var h hash.Hash

	switch strings.ToLower(algorithm) {
	case "sha256": h = sha256.New() // Most frequent case the first
	case "md5": h = md5.New()
	case "sha1": h = sha1.New()
	case "sha512": h = sha512.New()
	case "sha3-256": h = sha3.New256()
	case "sha3-512": h = sha3.New512()
	default:
		return nil, fmt.Errorf("hash algorithm %s is not supported", algorithm)
	}

	return &hasher{
		hash: h,
		buf: make([]byte, bufferSize),
	}, nil
}

func (h *hasher) hashFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file \"%s\": %s", path, err.Error())
	}
	defer f.Close()

	h.hash.Reset()
	if _, err := io.CopyBuffer(h.hash, f, h.buf); err != nil {
		return nil, fmt.Errorf("error hashing file \"%s\": %s", path, err.Error())
	}
	return h.hash.Sum(nil), nil
}

func (h *hasher) getFile(path string) (file, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return file{}, fmt.Errorf("cannot get information of \"%s\": %s", path, err.Error())
	}

	hash, err := h.hashFile(path)
	if err != nil {
		return file{}, err
	}

	return file{
		Name: filepath.Base(path),
		Size: stat.Size(),
		Hash: hash,
		realPath: path,
	}, nil
}

func (h *hasher) fileChecker(in *safeStringList, errs *safeCounter, wg sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for {
		path := in.next()
		if path == nil {
			break
		}

		stat, err := os.Stat(*path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot get info from \"%s\": %s\n", *path, err.Error())
			errs.increase()
			continue
		}

		f, err := getFileFromName(stat.Name())
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			errs.increase()
			continue
		}

		if f.Size != stat.Size() {
			fmt.Fprintf(os.Stderr, "sizes don't match in \"%s\"\n", *path)
			errs.increase()
			continue
		}

		hash, err := h.hashFile(*path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			errs.increase()
			continue
		}

		if !bytes.Equal(f.Hash, hash) {
			fmt.Fprintf(os.Stderr, "hashes don't match in \"%s\"\n", *path)
			errs.increase()
			continue
		} else {
			// TODO log OK when implement logger
		}
	}
}

func (mh *multiHasher) checkFiles(paths []string) (errs int) {
	var wg sync.WaitGroup
	var safeErrs safeCounter
	pathsSafe := makeConcurrentStringList(paths)

	for _, w := range mh.workers {
		go w.fileChecker(pathsSafe, &safeErrs, wg)
	}

	wg.Wait()
	return safeErrs.value()
}

func (h *hasher) fileGetter(in *safeStringList, out *safeFileList) error {
	for {
		path := in.next()
		if path == nil {
			break
		}

		f, err := h.getFile(*path)
		if err != nil {
			if OmitErrors {
				fmt.Fprintf(os.Stderr, "Error hashing file \"%s\": %s\n", *path, err.Error())
				continue
			} else {
				return err
			}
		}
		out.append(f)
	}

	return nil
}

func (mh *multiHasher) getFiles(paths []string) ([]file, error) {
	var eg errgroup.Group
	pathsSafe := makeConcurrentStringList(paths)
	filesSafe := makeConcurrentFileList(make([]file, 0, len(paths)))

	for _, w := range mh.workers {
		eg.Go(func() error {
			return w.fileGetter(pathsSafe, filesSafe)
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return filesSafe.getList(), nil
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

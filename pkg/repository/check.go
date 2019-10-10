package repository

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup/pkg/repository/settings"
	"github.com/Miguel-Dorta/gkup/pkg/threadSafe"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"hash"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func Check(path string, bufSize int) error {
	if bufSize < 512 {
		bufSize = 512
	}

	// Get settings (will be used later)
	sett, err := settings.Read(filepath.Join(path, settings.FileName))
	if err != nil {
		return fmt.Errorf("error reading settings: %w", err)
	}

	// Get all files
	fileList, err := getAllFiles(path)
	if err != nil {
		return fmt.Errorf("error listing repository files: %w", err)
	}
	safeFileList := threadSafe.NewStringList(fileList)

	// Do concurrent check
	errList := threadSafe.NewErrorList(nil)
	wg := &sync.WaitGroup{}
	for i:=0; i<runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := make([]byte, bufSize)
			h, err := getHash(sett.HashAlgorithm)
			if err != nil {
				errList.Append(err)
				return
			}

			for {
				f := safeFileList.Next()
				if f == nil {
					break
				}
				if err := checkFile(*f, h, buf); err != nil {
					errList.Append(err)
				}
			}
		}()
	}
	wg.Wait()

	printErrs(errList.GetList())
	return nil
}

func printErrs(errs []error) {
	if len(errs) == 0 {
		return
	}
	for _, err := range errs {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}

func getHash(algorithm string) (hash.Hash, error) {
	switch strings.ToLower(algorithm) {
	case "sha256":
		return sha256.New(), nil
	default:
		return nil, errors.New("hash algorithm unknown")
	}
}

func checkFile(path string, h hash.Hash, buf []byte) error {
	// Get data
	expectedHash, expectedSize, err := files.GetDataFromName(filepath.Base(path))
	if err != nil {
		return &os.PathError{
			Op:   "get data from filename",
			Path: path,
			Err:  err,
		}
	}

	// Check size
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot access file info: %w", err)
	}
	if stat.Size() != expectedSize {
		return fmt.Errorf("sizes don't match in file %s", path)
	}

	// Check hash
	actualHash, err := hashFile(path, h, buf)
	if err != nil {
		return fmt.Errorf("error hashing file: %w", err)
	}
	if !bytes.Equal(actualHash, expectedHash) {
		return fmt.Errorf("hashes don't match in file %s", path)
	}

	return nil
}

func getAllFiles(path string) ([]string, error) {
	result := make([]string, 0, 10000)

	filesFolderPath := filepath.Join(path, filesFolderName)
	for i:=0; i<=0xff; i++ {
		dirPath := filepath.Join(filesFolderPath, fmt.Sprintf("%02x", i))

		// List dir
		fList, err := utils.ListDir(dirPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, &os.PathError{
				Op:   "getAllRepoFiles",
				Path: dirPath,
				Err:  err,
			}
		}

		// Add files to list
		for _, f := range fList {
			if !f.Mode().IsRegular() {
				continue
			}
			result = append(result, filepath.Join(dirPath, f.Name()))
		}
	}
	return result, nil
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

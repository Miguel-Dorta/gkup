package hasher

import (
	"bytes"
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"github.com/Miguel-Dorta/gkup/pkg/logger"
	"github.com/Miguel-Dorta/gkup/pkg/threadSafe"
	"os"
	"sync"
)

// fileChecker is a worker that reads paths and checks whether the files listed in the path slice provided match with the info contained in their names.
// That means that they follow the specification from files.GetFileFromName() and their information is correct.
// This process is aimed to detect file corruption or filename defects.
// It returns the number of errors found
func (h *Hasher) fileChecker(in *threadSafe.StringList, errsFound *threadSafe.Fuse, wg sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for {
		path := in.Next()
		if path == nil {
			break
		}

		logger.Log.Debugf("Checking integrity of %s", *path)
		stat, err := os.Stat(*path)
		if err != nil {
			logger.Log.Errorf("cannot get info from \"%s\": %s\n", *path, err.Error())
			errsFound.Trigger()
			continue
		}

		f, err := files.GetFileFromName(stat.Name())
		if err != nil {
			logger.Log.Error(err.Error())
			errsFound.Trigger()
			continue
		}

		if f.Size != stat.Size() {
			logger.Log.Errorf("sizes don't match in \"%s\"\n", *path)
			errsFound.Trigger()
			continue
		}

		hash, err := h.HashPath(*path)
		if err != nil {
			logger.Log.Error(err.Error())
			errsFound.Trigger()
			continue
		}

		if !bytes.Equal(f.Hash, hash) {
			logger.Log.Errorf("hashes don't match in \"%s\"\n", *path)
			errsFound.Trigger()
			continue
		}

		logger.Log.Debugf("File %s is correct", *path)
	}
}

// fileGetter is a worker that reads paths, gets its files.File, and write those last ones in a list.
func (h *Hasher) fileGetter(in *threadSafe.StringList, out *threadSafe.FileList) error {
	for {
		path := in.Next()
		if path == nil {
			break
		}

		f, err := h.GetFile(*path)
		if err != nil {
			if logger.OmitErrors {
				logger.Log.Errorf("Error hashing file \"%s\": %s\n", *path, err.Error())
				continue
			} else {
				return err
			}
		}
		out.Append(f)
	}

	return nil
}

// fileHasher is a worker that gets and assigns the hash from the files.File provided
func (h *Hasher) fileHasher(list *threadSafe.FileList) error {
	for {
		f := list.Next()
		if f == nil {
			break
		}

		if err := h.HashFile(f); err != nil {
			if logger.OmitErrors {
				logger.Log.Error(err.Error())
				continue
			} else {
				return err
			}
		}
	}

	return nil
}

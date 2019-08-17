package hasher

import (
	"bytes"
	"github.com/Miguel-Dorta/gkup/pkg"
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"github.com/Miguel-Dorta/gkup/pkg/threadSafe"
	"os"
	"sync"
)

// fileChecker is a worker that reads paths and checks whether the files listed in the path slice provided match with the info contained in their names.
// That means that they follow the specification from files.GetFileFromName() and their information is correct.
// This process is aimed to detect file corruption or filename defects.
// It returns the number of errors found
func (h *Hasher) fileChecker(in *threadSafe.StringList, errsFound *fuse, wg *sync.WaitGroup) {
	for {
		path := in.Next()
		if path == nil {
			break
		}

		pkg.Log.Debugf("Checking integrity of %s", *path)
		stat, err := os.Stat(*path)
		if err != nil {
			pkg.Log.Errorf("cannot get info from \"%s\": %s\n", *path, err.Error())
			errsFound.trigger()
			continue
		}

		f, err := files.GetFileFromName(stat.Name())
		if err != nil {
			pkg.Log.Error(err.Error())
			errsFound.trigger()
			continue
		}

		if f.Size != stat.Size() {
			pkg.Log.Errorf("sizes don't match in \"%s\"\n", *path)
			errsFound.trigger()
			continue
		}

		hash, err := h.HashPath(*path)
		if err != nil {
			pkg.Log.Error(err.Error())
			errsFound.trigger()
			continue
		}

		if !bytes.Equal(f.Hash, hash) {
			pkg.Log.Errorf("hashes don't match in \"%s\"\n", *path)
			errsFound.trigger()
			continue
		}

		pkg.Log.Debugf("File %s is correct", *path)
	}
	wg.Done()
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
			if pkg.OmitErrors {
				pkg.Log.Errorf("Error hashing file \"%s\": %s\n", *path, err.Error())
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
			if pkg.OmitErrors {
				pkg.Log.Error(err.Error())
				continue
			} else {
				return err
			}
		}
	}

	return nil
}

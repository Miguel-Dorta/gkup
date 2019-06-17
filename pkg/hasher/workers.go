package hasher

import (
	"bytes"
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"github.com/Miguel-Dorta/gkup/pkg/threadSafe"
	"github.com/Miguel-Dorta/gkup/pkg/tmp"
	"os"
	"sync"
)

// fileChecker is a worker that reads paths and checks whether the files listed in the path slice provided match with the info contained in their names.
// That means that they follow the specification from files.GetFileFromName() and their information is correct.
// This process is aimed to detect file corruption or filename defects.
// It returns the number of errors found
func (h *Hasher) fileChecker(in *threadSafe.StringList, errs *threadSafe.Counter, wg sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for {
		path := in.Next()
		if path == nil {
			break
		}

		stat, err := os.Stat(*path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot get info from \"%s\": %s\n", *path, err.Error())
			errs.Increase()
			continue
		}

		f, err := files.GetFileFromName(stat.Name())
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			errs.Increase()
			continue
		}

		if f.Size != stat.Size() {
			fmt.Fprintf(os.Stderr, "sizes don't match in \"%s\"\n", *path)
			errs.Increase()
			continue
		}

		hash, err := h.HashPath(*path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			errs.Increase()
			continue
		}

		if !bytes.Equal(f.Hash, hash) {
			fmt.Fprintf(os.Stderr, "hashes don't match in \"%s\"\n", *path)
			errs.Increase()
			continue
		} else {
			// TODO log OK when implement logger
		}
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
			if tmp.OmitErrors {
				fmt.Fprintf(os.Stderr, "Error hashing file \"%s\": %s\n", *path, err.Error())
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
			if tmp.OmitErrors {
				os.Stderr.WriteString(err.Error() + "\n")
				continue
			} else {
				return err
			}
		}
	}

	return nil
}
package repo

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg/hasher"
	"github.com/Miguel-Dorta/gkup/pkg/tmp"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"os"
	"path/filepath"
	"runtime"
)

// CheckIntegrity checks the integrity of the files stored in the repo
func (r *Repo) CheckIntegrity() (errs int) {
	if r.sett == nil {
		fmt.Fprintln(os.Stderr, "Error: settings not loaded")
		return 1
	}

	// l1 is the list of elements in repo/files.
	// It should contain folders named from 00 to ff.
	l1, err := utils.ListDir(r.filesFolder)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error accessing files. Aborting...")
		return 1
	}

	mHasher, err := hasher.NewMultiHasher(r.sett.HashAlgorithm, tmp.BufferSize, runtime.NumCPU())
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return 1
	}

	allFiles := make([]string, 0, 1000)

	// c1 represent a given Child of the list l1
	for _, c1 := range l1 {
		if !c1.IsDir() {
			continue
		}
		c1Path := filepath.Join(r.filesFolder, c1.Name())

		// l2 is the list of elements in repo/files/c1.
		// It should contain the files named like <hash_hex>-<size_bytes>
		l2, err := utils.ListDir(c1Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error listing \"%s\": %s\n", c1Path, err.Error())
			errs++
			continue
		}
		// c2 represents a given Child of the list l2,
		// that means, a file in repo/files/c1
		for _, c2 := range l2 {
			if !c2.Mode().IsRegular() {
				continue
			}
			allFiles = append(allFiles, filepath.Join(c1Path, c2.Name()))
		}
	}

	errs += mHasher.CheckFiles(allFiles)
	return errs
}


package repo

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg/hasher"
	"github.com/Miguel-Dorta/gkup/pkg/logger"
	"github.com/Miguel-Dorta/gkup/pkg/tmp"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"path/filepath"
	"runtime"
)

// CheckIntegrity checks the integrity of the files stored in the repo
func (r *Repo) CheckIntegrity() error {
	if r.sett == nil {
		return errors.New("settings not loaded")
	}

	// l1 is the list of elements in repo/files.
	// It should contain folders named from 00 to ff.
	l1, err := utils.ListDir(r.filesFolder)
	if err != nil {
		return errors.New("cannot access files")
	}

	mHasher, err := hasher.NewMultiHasher(r.sett.HashAlgorithm, tmp.BufferSize, runtime.NumCPU())
	if err != nil {
		return err
	}

	allFiles := make([]string, 0, 1000)

	errsFound := false
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
			logger.Log.Error(fmt.Sprintf("error listing \"%s\": %s\n", c1Path, err.Error()))
			errsFound = true
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

	if mHasher.CheckFiles(allFiles) || errsFound {
		return errors.New("some errors were found")
	}
	return nil
}


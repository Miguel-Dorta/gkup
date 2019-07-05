package repo

import (
	"errors"
	"github.com/Miguel-Dorta/gkup/pkg"
	"github.com/Miguel-Dorta/gkup/pkg/hasher"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"path/filepath"
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

	mHasher, err := hasher.NewMultiHasher(r.sett.HashAlgorithm)
	if err != nil {
		return err
	}

	allFiles := make([]string, 0, 1000)

	pkg.Log.Info("Listing files")
	errsFound := false
	// c1 represent a given Child of the list l1
	for _, c1 := range l1 {
		c1Path := filepath.Join(r.filesFolder, c1.Name())
		if !c1.IsDir() {
			pkg.Log.Debugf("%s is not a directory. Skipping...", c1Path)
			continue
		}

		// l2 is the list of elements in repo/files/c1.
		// It should contain the files named like <hash_hex>-<size_bytes>
		l2, err := utils.ListDir(c1Path)
		if err != nil {
			pkg.Log.Errorf("error listing \"%s\": %s\n", c1Path, err.Error())
			errsFound = true
			continue
		}
		// c2 represents a given Child of the list l2,
		// that means, a file in repo/files/c1
		for _, c2 := range l2 {
			c2Path := filepath.Join(c1Path, c2.Name())
			if !c2.Mode().IsRegular() {
				pkg.Log.Debugf("%s is not a file. Skipping...", c2Path)
				continue
			}
			allFiles = append(allFiles, c2Path)
		}
	}

	pkg.Log.Info("Checking file integrity")
	if mHasher.CheckFiles(allFiles) || errsFound {
		return errors.New("some errors were found")
	}
	return nil
}


package repository

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg"
	"github.com/Miguel-Dorta/gkup/pkg/repository/settings"
	"os"
	"path/filepath"
)

const (
	backupsFolderName = "backups"
	filesFolderName = "files"
)

func Create(path, hashAlgorithm string) error {
	// Get path stat
	stat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) { // If it's not a "not exist" error, return it
			return &os.PathError{
				Op:   "stat repository path",
				Path: path,
				Err:  err,
			}
		}
		return create(path, hashAlgorithm) // If it's a "not exist" error, create it. Done.
	}

	// Check if it's not a dir
	if !stat.IsDir() {
		return &os.PathError{
			Op:   "create repository",
			Path: path,
			Err:  errors.New("must be a directory"),
		}
	}

	return create(path, hashAlgorithm)
}

// create creates a repository in the path provided with the algorithm provided.
// the path must exist and be an empty directory.
func create(path, hashAlgorithm string) error {
	// Create backups dir
	backupsFolderPath := filepath.Join(path, backupsFolderName)
	if err := os.MkdirAll(backupsFolderPath, pkg.DefaultDirPerm); err != nil {
		return &os.PathError{
			Op:   "create backup folder",
			Path: backupsFolderPath,
			Err:  err,
		}
	}

	// Create files dir and subdirectories
	filesFolderPath := filepath.Join(path, filesFolderName)
	for i:=0; i<=0xff; i++ {
		subDirPath := filepath.Join(filesFolderPath, fmt.Sprintf("%02x", i))
		if err := os.MkdirAll(subDirPath, pkg.DefaultDirPerm); err != nil {
			return &os.PathError{
				Op:   "create files folders",
				Path: subDirPath,
				Err:  err,
			}
		}
	}

	// Create settings file
	if err := settings.Write(filepath.Join(path, settings.FileName), hashAlgorithm); err != nil {
		return fmt.Errorf("error creating settings file: %s", err)
	}
	return nil
}

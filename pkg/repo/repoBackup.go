package repo

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"github.com/Miguel-Dorta/gkup/pkg/hasher"
	"github.com/Miguel-Dorta/gkup/pkg/logger"
	"github.com/Miguel-Dorta/gkup/pkg/tmp"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// BackupPaths backs up the paths provided and save the backup info in a file in BackupFolderName with the moment where it was created as name.
func (r *Repo) BackupPaths(paths []string) error {
	if r.sett == nil {
		return errors.New("settings not loaded")
	}

	now := time.Now() //Save the moment where the backup started

	fileList := make([]*files.File, 0, 1000)
	b := backup{
		Files: make([]*files.File, 0, 10),
		Dirs: make([]files.Dir, 0, 10),
	}

	logger.Log.Info("Listing files")
	// List all paths recursively
	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("\"%s\" not found", path)
			}
			return err
		}

		if stat.Mode().IsDir() {
			logger.Log.Debugf("Listing directory %s", path)
			child, childFiles, err := files.NewDir(path)
			if err != nil {
				if tmp.OmitErrors {
					logger.Log.Error(err.Error())
					continue
				} else {
					return err
				}
			}
			b.Dirs = append(b.Dirs, child)
			fileList = append(fileList, childFiles...)
		} else if stat.Mode().IsRegular() {
			logger.Log.Debugf("Listing file %s", path)
			child, err := files.NewFile(path)
			if err != nil {
				if tmp.OmitErrors {
					logger.Log.Error(err.Error())
					continue
				} else {
					return err
				}
			}
			b.Files = append(b.Files, child)
		} else {
			// TODO symlinks and other things
		}
	}
	fileList = append(fileList, b.Files...)

	// Get hash from all files
	multiH, err := hasher.NewMultiHasher(r.sett.HashAlgorithm, tmp.BufferSize, runtime.NumCPU())
	if err != nil {
		return err
	}
	logger.Log.Info("Hashing files")
	if err = multiH.HashFiles(fileList); err != nil {
		return err
	}

	logger.Log.Info("Adding files to repo")
	// Copy all files to repo
	for _, f := range fileList {
		if err := r.addFile(f); err != nil {
			if tmp.OmitErrors {
				logger.Log.Error(err.Error())
				continue
			} else {
				return err
			}
		}
	}

	logger.Log.Info("Saving backup")
	if err = writeBackup(
		filepath.Join(
			r.backupFolder,
			fmt.Sprintf(
				"%04d-%02d-%02d_%02d-%02d-%02d.json",
				now.Year(),
				now.Month(),
				now.Day(),
				now.Hour(),
				now.Minute(),
				now.Second(),
			),
		),
	b); err != nil {
		return err
	}

	return nil
}

// addFile adds a file to the file store of the repo
func (r *Repo) addFile(f *files.File) error {
	logger.Log.Debugf("Adding file %s to repo", f.RealPath)
	pathToSave := r.getPathInRepo(f)

	// If file already exists, do nothing. If exists but there's an error, return it
	if _, err := os.Stat(pathToSave); err == nil {
		logger.Log.Debug("It's already in the repo. Omitting...")
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("cannot get information of \"%s\": %s", pathToSave, err.Error())
	}

	if err := utils.CopyFile(f.RealPath, pathToSave); err != nil {
		return err
	}
	return nil
}

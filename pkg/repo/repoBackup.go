package repo

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup/pkg"
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"github.com/Miguel-Dorta/gkup/pkg/hasher"
	"github.com/Miguel-Dorta/gkup/pkg/utils"
	"os"
	"path/filepath"
	"time"
)

// BackupPaths backs up the paths provided and save the backup info in a file in BackupFolderName with the moment where it was created as name.
func (r *Repo) BackupPaths(paths []string, backupName string) error {
	if r.sett == nil {
		panic("settings not loaded in BackupPaths")
	}

	startTime := time.Now() //Save the moment where the backup started

	fileList := make([]*files.File, 0, pkg.SliceBigCapacity)
	b := backupFile{
		Files: make([]*files.File, 0, pkg.SliceSmallCapacity),
		Dirs:  make([]files.Dir, 0, pkg.SliceSmallCapacity),
	}

	// List all paths recursively
	pkg.Log.Info("Listing files")
	for _, path := range paths {

		// Get info of file (and check if it exists)
		stat, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("\"%s\" not found", path)
			}
			return err
		}

		// Skip if it's hidden
		if pkg.OmitHidden && utils.IsHidden(filepath.Base(path)) {
			pkg.Log.Debugf("omitting hidden file %s", path)
			continue
		}

		if stat.Mode().IsDir() {
			pkg.Log.Debugf("Listing directory %s", path)
			child, childFiles, err := files.NewDir(path)
			if err != nil {
				if pkg.OmitErrors {
					pkg.Log.Error(err.Error())
					continue
				} else {
					return err
				}
			}

			b.Dirs = append(b.Dirs, child)
			fileList = append(fileList, childFiles...)

		} else if stat.Mode().IsRegular() {
			pkg.Log.Debugf("Listing file %s", path)
			child, err := files.NewFile(path)
			if err != nil {
				if pkg.OmitErrors {
					pkg.Log.Error(err.Error())
					continue
				} else {
					return err
				}
			}

			b.Files = append(b.Files, child)

		} else {
			pkg.Log.Debugf("omitting unsupported file %s", path)
		}
	}
	fileList = append(fileList, b.Files...)

	// Get hash from all files
	multiH, err := hasher.NewMultiHasher(r.sett.HashAlgorithm)
	if err != nil {
		return err
	}
	pkg.Log.Info("Hashing files")
	if err = multiH.HashFiles(fileList); err != nil {
		return err
	}

	pkg.Log.Info("Adding files to repo")
	copyBuffer := make([]byte, pkg.BufferSize)
	// Copy all files to repo
	for _, f := range fileList {
		if err := r.addFile(f, copyBuffer); err != nil {
			if pkg.OmitErrors {
				pkg.Log.Error(err.Error())
				continue
			} else {
				return err
			}
		}
	}

	backupFileName := fmt.Sprintf(
		"%04d-%02d-%02d_%02d-%02d-%02d.json",
		startTime.Year(),
		startTime.Month(),
		startTime.Day(),
		startTime.Hour(),
		startTime.Minute(),
		startTime.Second(),
	)

	if backupName != "" {
		backupFileName = filepath.Join(backupName, backupFileName)
	}

	pkg.Log.Info("Saving backup")
	if err = writeBackup(filepath.Join(r.backupFolder, backupFileName), b); err != nil {
		return err
	}

	return nil
}

// addFile adds a file to the file store of the repo
func (r *Repo) addFile(f *files.File, buffer []byte) error {
	pkg.Log.Debugf("Adding file %s to repo", f.RealPath)
	pathToSave := r.getPathInRepo(f)

	// If file already exists, do nothing. If exists but there's an error, return it
	if _, err := os.Stat(pathToSave); err == nil {
		pkg.Log.Debug("It's already in the repo. Omitting...")
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("cannot get information of \"%s\": %s", pathToSave, err.Error())
	}

	if err := utils.CopyFile(f.RealPath, pathToSave, buffer); err != nil {
		return err
	}
	return nil
}

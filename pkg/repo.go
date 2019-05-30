package pkg

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func CreateRepo() (err error) {
	// Check if it's a directory
	stat, err := os.Stat(RepoPath)
	if err != nil {
		return
	}
	if !stat.IsDir() {
		return //TODO not dir
	}

	// Check if it's empty
	childs, err := listDir(RepoPath)
	if err != nil {
		return
	}
	if len(childs) != 0 {
		return //TODO not empty
	}

	// Make backup folder
	err = os.Mkdir(filepath.Join(RepoPath, "backups"), 0700)
	if err != nil {
		return
	}

	// Make file folder and subforders
	fileDirPath := filepath.Join(RepoPath, "files")
	err = os.Mkdir(fileDirPath, 0700)
	if err != nil {
		return
	}
	for i:=0x0; i<=0xff; i++ {
		err = os.Mkdir(filepath.Join(fileDirPath, fmt.Sprintf("%02x", i)), 0700)
		if err != nil {
			return
		}
	}

	// Write settings.toml
	data, err := generateSettingsToml(false)
	if err != nil {
		return
	}
	return ioutil.WriteFile(filepath.Join(RepoPath, "settings.toml"), data, 0600)
}

func CheckIntegrity() []error {
	errs := make([]error, 0, 10)

	fileDirPath := filepath.Join(RepoPath, "files")
	l1, err := listDir(fileDirPath)
	if err != nil {
		return append(errs, err)
	}
	for _, c1 := range l1 {
		if !c1.IsDir() {
			continue
		}

		fileByteDirPath := filepath.Join(fileDirPath, c1.Name())
		l2, err := listDir(fileByteDirPath)
		if err != nil {
			errs = append(errs, err) //TODO
		}

		for _, c2 := range l2 {
			if c2.IsDir() {
				continue
			}
			name := c2.Name()

			// Find hash-size
			var hash []byte
			var size int64 = -1
			for i, b := range name {
				if b != '-' {
					continue
				}

				hash, err = hex.DecodeString(name[:i])
				if err != nil {
					errs = append(errs, err)
					break //TODO
				}

				size, err = strconv.ParseInt(name[i+1:], 10, 64)
				if err != nil {
					errs = append(errs, err)
					break //TODO
				}
			}
			if len(hash) == 0 || size < 0 {
				continue //TODO
			}

			if c2.Size() != size {
				//TODO
			}

			newHash, err := hashFile(filepath.Join(fileByteDirPath, name))
			if err != nil {
				//TODO
			}

			if !equals(hash, newHash) {
				//TODO
			}
		}
	}
	return errs
}

func BackupPaths(paths []string) error {
	now := time.Now() //Save the moment where the backup started

	//TODO check for duplicates

	d := dir{
		Files: make([]file, 0, 10),
		Dirs: make([]dir, 0, 10),
	}
	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			return err //TODO check if exists
		}
		mode := stat.Mode()
		name := stat.Name()

		if mode.IsDir() {
			child, err := listFilesRecursive(path, name)
			if err != nil {
				return err
			}
			d.Dirs = append(d.Dirs, child)
		} else if mode.IsRegular() {
			child, err := getFile(path)
			if err != nil {
				return err
			}

			if err = addFile(child); err != nil {
				return err
			}

			d.Files = append(d.Files, child)
		} else {
			// TODO symlinks and other things
		}
	}

	return writeBackup(fmt.Sprintf(
		"%04d-%02d-%02d_%02d-%02d-%02d.json",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
	), d)
}

func RestoreBackup(date, restoreTo string) error {
	var b backup
	{
		backupPath := filepath.Join(RepoPath, "backups")
		backupList, err := listDir(backupPath)
		if err != nil {
			return err
		}

		found := false
		for _, bac := range backupList {
			name := bac.Name()
			if strings.HasPrefix(name, date) {
				backupPath = filepath.Join(backupPath, name)
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("") //TODO not found
		}

		b, err = readBackup(backupPath)
		if err != nil {
			return err
		}
	}

	//TODO check versioning

	return restoreDir(dir{Files: b.Files, Dirs: b.Dirs}, restoreTo)
}

func restoreDir(d dir, pathToRestore string) (err error) {
	for _, childFile := range d.Files {
		if err = copyFile(getPathInRepo(childFile), filepath.Join(pathToRestore, childFile.Name)); err != nil {
			return
		}
	}

	for _, childDir := range d.Dirs {
		childPath := filepath.Join(pathToRestore, childDir.Name)
		if err = os.Mkdir(childPath, 0700); err != nil {
			return
		}
		if err = restoreDir(childDir, childPath); err != nil {
			return
		}
	}

	return
}

func getPathInRepo(f file) string {
	hashStr := hex.EncodeToString(f.Hash)
	return filepath.Join(
		RepoPath,
		"files",
		hashStr[:2],
		fmt.Sprintf("%s-%d", hashStr, f.Size),
	)
}

func addFile(f file) error {
	pathToSave := getPathInRepo(f)

	// If file already exists, do nothing. If exists but there's an error, return it
	if _, err := os.Stat(pathToSave); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	err := copyFile(f.realPath, pathToSave)
	if err != nil {
		return err //TODO critical failure. If fails during copy, there'll be a ghost file
	}
	return nil
}

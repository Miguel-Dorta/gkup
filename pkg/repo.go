package pkg

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	backupFolder = "backup"
	filesFolder = "files"
	settingsPath = "settings.toml"
)

func CreateRepo() error {
	// Check if it's a directory
	if stat, err := os.Stat(RepoPath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error getting stats from \"%s\": %s", RepoPath, err.Error())
		}
		// If not exists, create it
		if err = os.MkdirAll(RepoPath, 0700); err != nil {
			return fmt.Errorf("error creating directory \"%s\": %s", RepoPath, err.Error())
		}
	} else if !stat.IsDir() {
		//TODO read symlink
		return fmt.Errorf("\"%s\" is not a directory", RepoPath)
	}

	// Check if it's empty
	if children, err := listDir(RepoPath); err != nil {
		return fmt.Errorf("error listing \"%s\": %s", RepoPath, err.Error())
	} else if len(children) != 0 {
		return fmt.Errorf("\"%s\" is not empty", RepoPath)
	}

	// Make backup folder
	if err := os.Mkdir(backupFolder, 0700); err != nil {
		return fmt.Errorf("error creating subdirectory \"%s\": %s", backupFolder, err.Error())
	}

	// Make file folder and subforders
	if err := os.Mkdir(filesFolder, 0700); err != nil {
		return fmt.Errorf("error creating subdirectory \"%s\": %s", filesFolder, err.Error())
	}
	for i:=0x0; i<=0xff; i++ {
		path := filepath.Join(filesFolder, fmt.Sprintf("%02x", i))
		if err := os.Mkdir(path, 0700); err != nil {
			return fmt.Errorf("error creating subdirectory \"%s\": %s", path, err.Error())
		}
	}

	// Write settings.toml
	if err := ioutil.WriteFile(settingsPath, generateSettingsToml(), 0600); err != nil {
		return fmt.Errorf("error writing settings: %s", err.Error())
	}
	return nil
}

func CheckIntegrity() (errs int) {
	//TODO load settings

	// l1 is the list of elements in repo/files.
	// It should contain folders named from 00 to ff.
	l1, err := listDir(filesFolder)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error accessing files. Aborting...")
		return 1
	}
	// c1 represent a given Child of the list l1
	for _, c1 := range l1 {
		if !c1.IsDir() {
			continue
		}
		c1Path := filepath.Join(filesFolder, c1.Name())

		// l2 is the list of elements in repo/files/c1.
		// It should contain the files named like <hash_hex>-<size_bytes>
		l2, err := listDir(c1Path)
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
			c2Name := c2.Name()
			c2Path := filepath.Join(c1Path, c2Name)

			hash, size, err := getHashSize(c2Name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error getting info from \"%s\": %s\n", c2Path, err.Error())
				errs++
				continue
			}

			if c2.Size() != size {
				fmt.Fprintf(os.Stderr, "sizes don't match in \"%s\"\n", c2Path)
				errs++
				continue
			}

			if newHash, err := hashFile(c2Path); err != nil {
				fmt.Fprintf(os.Stderr, "cannot hash \"%s\"\n", c2Path)
				errs++
				continue
			} else if bytes.Equal(hash, newHash) {
				fmt.Fprintf(os.Stderr, "hashes don't match in \"%s\"\n", c2Path)
				errs++
				continue
			}
		}
	}
	return errs
}

func getHashSize(fileName string) (hash []byte, size int64, err error) {
	size = -1
	for i, b := range fileName {
		if b != '-' {
			continue
		}

		if hash, err = hex.DecodeString(fileName[:i]); err != nil {
			return nil, 0, fmt.Errorf("cannot decode hash: %s", err.Error())
		}

		if size, err = strconv.ParseInt(fileName[i+1:], 10, 64); err != nil {
			return nil, 0, fmt.Errorf("cannot parse size: %s", err.Error())
		}
		break
	}

	if hash == nil || size < 0 {
		return nil, 0, errors.New("invalid format")
	}

	return
}

func BackupPaths(paths []string) error {
	now := time.Now() //Save the moment where the backup started

	//TODO check for duplicates

	root := dir{
		Files: make([]file, 0, 10),
		Dirs: make([]dir, 0, 10),
	}
	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("\"%s\" not found", path)
			}
			return err
		}
		mode := stat.Mode()
		name := stat.Name()

		if mode.IsDir() {
			child, err := listFilesRecursive(path)
			if err != nil {
				if OmitErrors {
					os.Stderr.WriteString(err.Error())
					continue
				} else {
					return err
				}
			}
			root.Dirs = append(root.Dirs, child)
			// TODO I'm not adding the children to the backup, am I?
		} else if mode.IsRegular() {
			child, err := getFile(path)
			if err != nil {
				if OmitErrors {
					os.Stderr.WriteString(err.Error())
					continue
				} else {
					return err
				}
			}

			if err = addFile(child); err != nil {
				return err
			}

			root.Files = append(root.Files, child)
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
	), root)
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
		filesFolder,
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
		return fmt.Errorf("cannot get information of \"%s\": %s", pathToSave, err.Error())
	}

	if err := copyFile(f.realPath, pathToSave); err != nil {
		return err
	}
	return nil
}

package pkg

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	backupFolderName = "backup"
	filesFolderName  = "files"
)

type Repo struct {
	sett *settings
	path, backupFolder, filesFolder, settingsPath string
}

func NewRepo(repoPath string) Repo {
	return Repo{
		path: repoPath,
		backupFolder: filepath.Join(repoPath, backupFolderName),
		filesFolder: filepath.Join(repoPath, filesFolderName),
		settingsPath: filepath.Join(repoPath, settingsFilename),
	}
}

func (r *Repo) LoadSettings() error {
	sett, err := loadSettings(r.settingsPath)
	if err != nil {
		return err
	}
	r.sett = &sett
	return nil
}

func (r *Repo) Create(hashAlgorithm string) error {
	// Check if it's a directory
	if stat, err := os.Stat(r.path); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error getting stats from \"%s\": %s", r.path, err.Error())
		}
		// If not exists, create it
		if err = os.MkdirAll(r.path, 0700); err != nil {
			return fmt.Errorf("error creating directory \"%s\": %s", r.path, err.Error())
		}
	} else if !stat.IsDir() {
		//TODO read symlink
		return fmt.Errorf("\"%s\" is not a directory", r.path)
	}

	// Check if it's empty
	if children, err := listDir(r.path); err != nil {
		return fmt.Errorf("error listing \"%s\": %s", r.path, err.Error())
	} else if len(children) != 0 {
		return fmt.Errorf("\"%s\" is not empty", r.path)
	}

	// Make backup folder
	if err := os.Mkdir(r.backupFolder, 0700); err != nil {
		return fmt.Errorf("error creating subdirectory \"%s\": %s", r.backupFolder, err.Error())
	}

	// Make file folder and subforders
	if err := os.Mkdir(r.filesFolder, 0700); err != nil {
		return fmt.Errorf("error creating subdirectory \"%s\": %s", r.filesFolder, err.Error())
	}
	for i:=0x0; i<=0xff; i++ {
		path := filepath.Join(r.filesFolder, fmt.Sprintf("%02x", i))
		if err := os.Mkdir(path, 0700); err != nil {
			return fmt.Errorf("error creating subdirectory \"%s\": %s", path, err.Error())
		}
	}

	// Write settings.toml
	if err := saveSettings(r.settingsPath, hashAlgorithm); err != nil {
		return err
	}
	return nil
}

func (r *Repo) CheckIntegrity() (errs int) {
	if r.sett == nil {
		fmt.Fprintln(os.Stderr, "Error: settings not loaded")
		return 1
	}

	// l1 is the list of elements in repo/files.
	// It should contain folders named from 00 to ff.
	l1, err := listDir(r.filesFolder)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error accessing files. Aborting...")
		return 1
	}
	// c1 represent a given Child of the list l1
	for _, c1 := range l1 {
		if !c1.IsDir() {
			continue
		}
		c1Path := filepath.Join(r.filesFolder, c1.Name())

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

			if newHash, err := hashFile(c2Path, r.sett.HashAlgorithm); err != nil {
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

func (r *Repo) BackupPaths(paths []string) error {
	if r.sett == nil {
		return errors.New("settings not loaded")
	}

	now := time.Now() //Save the moment where the backup started

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

		if stat.Mode().IsDir() {
			child, err := r.listFilesRecursive(path)
			if err != nil {
				if OmitErrors {
					os.Stderr.WriteString(err.Error() + "\n")
					continue
				} else {
					return err
				}
			}
			root.Dirs = append(root.Dirs, child)
		} else if stat.Mode().IsRegular() {
			child, err := r.getFile(path)
			if err != nil {
				if OmitErrors {
					os.Stderr.WriteString(err.Error() + "\n")
					continue
				} else {
					return err
				}
			}

			if err = r.addFile(child); err != nil {
				if OmitErrors {
					os.Stderr.WriteString(err.Error() + "\n")
					continue
				} else {
					return err
				}
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

func (r *Repo) RestoreBackup(date, restoreTo string) error {
	if r.sett == nil {
		return errors.New("settings not loaded")
	}

	var b backup
	{
		backupsList, err := listDir(r.backupFolder)
		if err != nil {
			return err
		}

		var backupPath string
		for _, bac := range backupsList {
			if strings.HasPrefix(bac.Name(), date) {
				backupPath = filepath.Join(r.backupFolder, bac.Name())
				break
			}
		}

		if backupPath == "" {
			return errors.New("backupPath not found")
		}

		if b, err = readBackup(backupPath); err != nil {
			return err
		}
	}

	//TODO check versioning

	if err := r.restoreDir(dir{Files: b.Files, Dirs: b.Dirs}, restoreTo); err != nil {
		return err
	}
	return nil
}

func (r *Repo) restoreDir(d dir, pathToRestore string) error {
	for _, childFile := range d.Files {
		if err := copyFile(r.getPathInRepo(childFile), filepath.Join(pathToRestore, childFile.Name)); err != nil {
			if OmitErrors {
				os.Stderr.WriteString(err.Error() + "\n")
				continue
			} else {
				return err
			}
		}
	}

	for _, childDir := range d.Dirs {
		childPath := filepath.Join(pathToRestore, childDir.Name)
		if err := os.Mkdir(childPath, 0700); err != nil {
			if OmitErrors {
				fmt.Fprintf(os.Stderr, "Error restoring folder \"%s\": %s\n", childPath, err.Error())
				continue
			} else {
				return err
			}
		}
		if err := r.restoreDir(childDir, childPath); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repo) getPathInRepo(f file) string {
	hashStr := hex.EncodeToString(f.Hash)
	return filepath.Join(
		r.filesFolder,
		hashStr[:2],
		fmt.Sprintf("%s-%d", hashStr, f.Size),
	)
}

func (r *Repo) addFile(f file) error {
	pathToSave := r.getPathInRepo(f)

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
